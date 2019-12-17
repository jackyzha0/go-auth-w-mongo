package db

import (
	"context"
	"errors"
	"log"
	"time"

	"../schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_uri = "mongodb://localhost:27017"

// Create DB Connection
var Ctx context.Context
var Client *mongo.Client

var ErrDocumentNotFound = errors.New("no documents found matching that filter")

func init() {
	// Initialize DB
	log.Printf("Attempting to connect to %q", DB_uri)
	conn, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clt, err := mongo.Connect(Ctx, options.Client().ApplyURI(DB_uri))

	if err != nil {
		log.Panic(err)
	}
	log.Print("Connection Established")

	// Change Package level vars
	Ctx = conn
	Client = clt
}

func FindOne(filter bson.M, db, clct string) (schemas.User, error) {
	var result schemas.User

	// Set DB and collection
	collection := Client.Database(db).Collection(clct)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Actual find operation
	errF := collection.FindOne(Ctx, filter).Decode(&result)

	// Check if demarshalled into struct properly
	if errF != nil {
		return result, ErrDocumentNotFound
	}

	return result, nil
}

// func find(filter bson.D, db, clct string) (bson.D, error) {
// 	// Set DB and collection
// 	collection := Client.Database(db).Collection(clct)
// 	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
//
//   cursor, err := collection.Find(Ctx, bson.{})
//
// 	return cursor, err
// }
