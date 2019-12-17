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

var ErrDocumentNotFound = errors.New("no found under that filter")
var ErrCouldNotUnMarshall = errors.New("could not find unmarshall data. try looking at your query?")
var ErrCursorIterationFailed = errors.New("an error occurred when iterating through a cursor")
var ErrCollectionNotDefined = errors.New("collection not defined in db.go")

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

func FindOne(filter bson.M, db, clct string) (res interface{}, err error) {
	// Set DB and collection
	collection := Client.Database(db).Collection(clct)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	switch clct {
	case "users":
		res, err = FindOneUser(filter, collection)
	}

	// failed to demarshal
	if err != nil {
		log.Fatalf("could not unmarshall document with filter %q", filter)
		return res, ErrCouldNotUnMarshall
	}

	return res, err
}

func FindOneUser(filter bson.M, collection *mongo.Collection) (schemas.User, error) {

	// Find Operation
	var result schemas.User
	err := collection.FindOne(Ctx, filter).Decode(&result)

	return result, err
}

func FindUsers(filter bson.D, db, clct string) ([]*schemas.User, error) {

	// Set DB and collection
	collection := Client.Database(db).Collection(clct)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Find Operation
	var result []*schemas.User
	cursor, err := collection.Find(Ctx, filter)

	if err != nil {
		return result, ErrDocumentNotFound
	}

	for cursor.Next(Ctx) {
		var user schemas.User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatalf("could not unmarshall document with filter %q", filter)
			return result, ErrCouldNotUnMarshall
		}
		result = append(result, &user)
	}

	if cursor.Err() != nil {
		log.Fatalf("cursor iteration failed with filter %q", filter)
		return result, ErrCursorIterationFailed
	}

	// Close cursor after we're done with it
	cursor.Close(Ctx)

	return result, nil
}
