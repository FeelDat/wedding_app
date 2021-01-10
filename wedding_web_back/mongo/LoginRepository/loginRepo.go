package LoginRepository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"wedding_project/internal/models"

	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRepository struct {
	log                   *zap.Logger
	LoginsCollection *mongo.Collection
}

func NewLoginRepository(log *zap.Logger, Db *mongo.Collection) *LoginRepository {

	return &LoginRepository{
		log:              log,
		LoginsCollection: Db,
	}
}

func (dr *LoginRepository) CheckUser(login, pass string) (*models.Login, error) {
	var result *models.Login
	err := dr.LoginsCollection.FindOne(context.TODO(), bson.M{"login": login}).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Password == pass {
		return result, nil
	}
	return nil, fmt.Errorf("password is incorrect")
}
