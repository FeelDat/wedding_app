package GuestRepository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"wedding_project/internal/models"

	//"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type GusetsRepository struct {
	log             *zap.Logger
	GuestCollection *mongo.Collection
}

func NewGuestsRepository(log *zap.Logger, Db *mongo.Collection) *GusetsRepository {

	return &GusetsRepository{
		log:             log,
		GuestCollection: Db,
	}
}

func (mc *GusetsRepository) GetListOfAllGuests() []*models.Guest {
	var results []*models.Guest
	cur, err := mc.GuestCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Guest
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	//fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results
}

func (mc *GusetsRepository) GetGuest(id string) *models.Guest {
	var result *models.Guest
	err := mc.GuestCollection.FindOne(context.TODO(), id).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (mc *GusetsRepository) CreateGuest(guestName, guestNumber string) {
	newGuest := make(map[string]string)
	newGuest["name"] = guestName
	newGuest["number"] = guestNumber
	newGuest["disposition"] = "0"
	insertResult, err := mc.GuestCollection.InsertOne(context.TODO(), newGuest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func (mc *GusetsRepository) UpdateGuest(id, guestName, guestNumber, disposition string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.M{
		"$set": bson.M{
			"name":   guestName,
			"number": guestNumber,
			"disposition": disposition,
		},
	}

	updateResult, err := mc.GuestCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func (mc *GusetsRepository) DeleteGuest(id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	deleteResult, err := mc.GuestCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func (mc *GusetsRepository) DropDisposition() {
	allGuests := mc.GetListOfAllGuests()
	for _, v := range allGuests {
		if v.Disposition != "0" {
			mc.UpdateGuest(v.Id.Hex(), v.Name, v.Number, "0")
		}
	}
	//fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}
