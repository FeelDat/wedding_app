package LoginmetaRepository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
	"time"
	"wedding_project/internal/models"
)

type LoginmetaRepository struct {
	log                   *zap.Logger
	LoginmetaCollection *mongo.Collection
}

func NewLoginMetaRepository(log *zap.Logger, Db *mongo.Collection) *LoginmetaRepository {

	return &LoginmetaRepository{
		log:              log,
		LoginmetaCollection: Db,
	}
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    time.Time
	RtExpires    time.Time
}

func (dr *LoginmetaRepository) ExpireSet(userid string, td *TokenDetails) error {
	var result *models.Login
	err := dr.LoginmetaCollection.FindOne(context.TODO(), bson.M{"uuid": userid}).Decode(&result)
	if err != nil {
		dr.log.Error(err.Error())
	} else {
		err = dr.DeleteAuth(userid)
		if err != nil {
			dr.log.Error(err.Error())
		}
	}
	at := td.AtExpires
	rt := td.RtExpires
	res, err := dr.LoginmetaCollection.InsertMany(context.TODO(), []interface{}{
		bson.D{
			{"expiredAt", bsonx.Time(at)},
			{"type", "access"},
			{"uuid", userid},
			{"token", td.AccessToken},
		},
		bson.D{
			{"expiredAt", bsonx.Time(rt)},
			{"type", "refresh"},
			{"uuid", userid},
			{"token", td.RefreshToken},
		},
    })
	if err != nil {
		return err
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", len(res.InsertedIDs))
	return nil
}

func (dr *LoginmetaRepository) FetchAuth(uuid string) (error) {
	//var results []*models.Login
	filterCursor, err := dr.LoginmetaCollection.Find(context.TODO(), bson.M{"uuid": uuid})
	if err != nil {
		return err
	}
	var episodesFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &episodesFiltered); err != nil {
		return err
	}
	for _, v := range episodesFiltered {
		if v["type"] == "access" {
			return nil
		}
	}
	return fmt.Errorf("do not have access token")
}

func (dr *LoginmetaRepository) DeleteAuth(userid string) error {
	filter := bson.M{"uuid": bson.M{"$eq": userid}}
	deleteResult, err := dr.LoginmetaCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}