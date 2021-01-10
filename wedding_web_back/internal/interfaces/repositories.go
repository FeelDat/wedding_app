package interfaces

import (
	"wedding_project/internal/models"
	LoginmetaRepository2 "wedding_project/mongo/LoginmetaRepository"
)

type GuestRepositry interface {
	GetListOfAllGuests() []*models.Guest
	GetGuest(id string) *models.Guest
	CreateGuest(guestName, guestNumber string)
	UpdateGuest(id, guestName, guestNumber, disposition string)
	DeleteGuest(id string)
	DropDisposition()
}

type LoginRepository interface {
	CheckUser(login, pass string) (*models.Login, error)
}

type LoginmetaRepository interface {
	ExpireSet(userid string, td *LoginmetaRepository2.TokenDetails) error
	DeleteAuth(userid string) error
	FetchAuth(userid string) error
}