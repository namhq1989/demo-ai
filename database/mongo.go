package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string) *mongo.Client {
	var (
		ctx = context.Background()
	)

	// use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// send a ping to confirm a successful connection
	if err = client.Database("admin").RunCommand(ctx, bson.M{"ping": 1}).Err(); err != nil {
		panic(err)
	}

	fmt.Printf("⚡️ [mongodb]: connected \n")

	return client
}

func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func ObjectIDFromString(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func IsValidObjectID(id string) bool {
	_, err := ObjectIDFromString(id)
	return err == nil
}

func SetDefaultPageLimit(page *int64, limit *int64) {
	if *page < 0 {
		*page = 0
	}

	if *limit < 0 || *limit > 50 {
		*limit = 20
	}
}

func ColHistory(db *mongo.Database) *mongo.Collection {
	return db.Collection("histories")
}
