package server

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"wedding_project/config"
	"wedding_project/internal/interfaces"
	"wedding_project/mongo/LoginmetaRepository"
)

type RestApiServer struct {
	r               *httprouter.Router
	log             *zap.Logger
	config          *config.Configuration
	GuestRepo       interfaces.GuestRepositry
	LoginRepo interfaces.LoginRepository
	LoginmetaRepo interfaces.LoginmetaRepository
}

func NewRestApiServer(logger *zap.Logger, config *config.Configuration, guestsRepo interfaces.GuestRepositry, loginRepo interfaces.LoginRepository, loginmetaRepo interfaces.LoginmetaRepository) *RestApiServer {
	return &RestApiServer{httprouter.New(), logger, config, guestsRepo, loginRepo, loginmetaRepo}
}

func (s *RestApiServer) ListenAndServe(addr string) error {
	s.routes()
	return http.ListenAndServe(addr, s)
}

func (s *RestApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	s.r.ServeHTTP(w, r)
}

func (s *RestApiServer) routes() {
	s.r.GET("/guests", s.TokenAuthMiddleware(s.getListOfGuests))
	s.r.GET("/guests/:id", s.TokenAuthMiddleware(s.getGuestInfo))
	s.r.POST("/guests", s.TokenAuthMiddleware(s.createNewGuest))
	s.r.PUT("/guests/:id", s.TokenAuthMiddleware(s.updateGuestInfo))
	s.r.DELETE("/guests/:id", s.TokenAuthMiddleware(s.deleteGuest))
	s.r.GET("/disposition", s.TokenAuthMiddleware(s.getListOfGuests))
	//s.r.PUT("/disposition/:id", s.updateDisposition)
	s.r.DELETE("/disposition", s.TokenAuthMiddleware(s.dropDisposition))
	s.r.POST("/login/:name/:password", s.loginGuest)
	s.r.POST("/logout", s.TokenAuthMiddleware(s.Logout))
	s.r.POST("/refresh", s.Refresh)
}

/*func (s *RestApiServer) statsMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {

		if value, ok := r.URL.Query()["beginTime"]; ok {
			vars = append(vars, httprouter.Param{Key: "beginTime", Value: value[0]})
		}
		if value, ok := r.URL.Query()["period"]; ok {
			vars = append(vars, httprouter.Param{Key: "period", Value: value[0]})
		}

		h(w, r, vars)
	}
}*/
func (s *RestApiServer) TokenAuthMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
		err := TokenValid(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r, vars)
	}
}

type User struct {
	ID uint64            `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *RestApiServer) loginGuest(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	name := vars.ByName("name")
	pass := vars.ByName("password")
	user, err := s.LoginRepo.CheckUser(name, pass)

	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ts, err := s.CreateToken(user.Id.Hex())
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	saveErr := s.LoginmetaRepo.ExpireSet(user.Id.Hex(), ts)
	if saveErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	tokens := map[string]string{
		"id": user.Id.Hex(),
		"token":  ts.AccessToken,
		"refToken": ts.RefreshToken,
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(tokens)
}

func (s *RestApiServer) CreateToken(userid string) (*LoginmetaRepository.TokenDetails, error) {
	td := &LoginmetaRepository.TokenDetails{}
	td.AtExpires = time.Now().Add(15*time.Minute)
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7)
	td.RefreshUuid = uuid.NewV4().String()

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(s.config.AccessSecret))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(s.config.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (string, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, ok := claims["user_id"].(string)
		//userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if !ok {
			return "", fmt.Errorf("user_id invalid")
		}
		return userId, nil
	}
	return "", err
}

func (s *RestApiServer) Refresh(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	/*body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	refreshToken := map[string]string{}
	err = json.Unmarshal(body, &refreshToken)
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	//verify the token
	token, err := jwt.Parse(refreshToken["refreshToken"], func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.RefreshSecret), nil
	})*/
	token, err := VerifyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		//Delete the previous Refresh Token
		delErr := s.LoginmetaRepo.DeleteAuth(refreshUuid)
		if delErr != nil { //if any goes wrong
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := s.CreateToken(userId)
		if  createErr != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//save the tokens metadata to redis
		saveErr := s.LoginmetaRepo.ExpireSet(userId, ts)
		if saveErr != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (s *RestApiServer) Logout(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	au, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	delErr := s.LoginmetaRepo.DeleteAuth(au)
	if delErr != nil { //if any goes wrong
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *RestApiServer) getListOfGuests(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	guests := s.GuestRepo.GetListOfAllGuests()
	jsonOutput, err := json.Marshal(guests)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info("successfully decoded guests")
	}
	w.Write(jsonOutput)
}

func (s *RestApiServer) getGuestInfo(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	guestId := vars.ByName("id")
	guest := s.GuestRepo.GetGuest(guestId)
	jsonOutput, err := json.Marshal(guest)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info("error while decoding guest")
	}
	w.Write(jsonOutput)
}

func (s *RestApiServer) createNewGuest(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = s.LoginmetaRepo.FetchAuth(tokenAuth)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	var guest map[string]string
	err = json.Unmarshal(body, &guest)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	s.GuestRepo.CreateGuest(guest["name"], guest["number"])
	//s.DispositionRepo.SaveDisposition()
}

func (s *RestApiServer) updateGuestInfo(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = s.LoginmetaRepo.FetchAuth(tokenAuth)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	guestId := vars.ByName("id")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	var guest map[string]string
	err = json.Unmarshal(body, &guest)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
 	s.GuestRepo.UpdateGuest(guestId, guest["name"], guest["number"], guest["disposition"])
}

func (s *RestApiServer) deleteGuest(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	guestId := vars.ByName("id")
	s.GuestRepo.DeleteGuest(guestId)
	//s.DispositionRepo.DeleteGuestDisposition(guestId)
}

func (s *RestApiServer) updateDisposition(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	guestId := vars.ByName("id")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	var guest map[string]string
	err = json.Unmarshal(body, &guest)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	s.GuestRepo.UpdateGuest(guestId, "", "", guest["disposition"])
}

func (s *RestApiServer) dropDisposition(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	//drop disposition (all guests in table 0)
	s.GuestRepo.DropDisposition()
}

/*func (s *RestApiServer) getDisposition(w http.ResponseWriter, r *http.Request, vars httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	var guests map[string][]string
	err = json.Unmarshal(body, &guests)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	s.DispositionRepo.GetDisposition()
}*/
