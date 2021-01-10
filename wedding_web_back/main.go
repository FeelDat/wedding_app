package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"wedding_project/config"
	"wedding_project/internal/server"
	"wedding_project/mongo/GuestRepository"
	"wedding_project/mongo/LoginRepository"
	"wedding_project/mongo/LoginmetaRepository"
)

func main() {
	logger := config.GetLogger()
	cfg := config.GetConfig(logger)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	mongoDb := client.Database("guestsDB")
	guestsRepo := GuestRepository.NewGuestsRepository(logger, mongoDb.Collection("guests"))
	loginRepo := LoginRepository.NewLoginRepository(logger, mongoDb.Collection("logins"))
	loginmetaRepo := LoginmetaRepository.NewLoginMetaRepository(logger, mongoDb.Collection("loginmeta"))
	serv := server.NewRestApiServer(logger, cfg, guestsRepo, loginRepo, loginmetaRepo)

	/*err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")*/
	log.Fatal(serv.ListenAndServe(cfg.Server.Address + ":" + strconv.Itoa(cfg.Server.Port)))
}
