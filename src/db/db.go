package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB             *mongo.Database
	SessionsCol    *mongo.Collection
	UsersCol       *mongo.Collection
	OTPSessionsCol *mongo.Collection
)

func ConnectToDb() *mongo.Client {
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		log.Fatal("Required env var 'MONGODB_URI' not found!")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the database")

	DB = client.Database("finetrack")
	SessionsCol = DB.Collection("sessions")
	UsersCol = DB.Collection("users")
	OTPSessionsCol = DB.Collection("otpsessions")

	// prepare the indexes
	userEmailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = UsersCol.Indexes().CreateOne(context.TODO(), userEmailIndex)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
