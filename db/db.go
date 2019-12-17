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

func FindOneUser(filter bson.M, db, clct string) (schemas.User, error) {
	var result schemas.User

	// Set DB and collection
	collection := Client.Database(db).Collection(clct)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Actual find operation
	err := collection.FindOne(Ctx, filter).Decode(&result)

	// failed to demarshal
	if err != nil {
		log.Fatalf("could not unmarshall document with filter %q", filter)
		return result, ErrCouldNotUnMarshall
	}

	return result, nil
}

func UpdateOneUser() {

}

func FindUsers(filter bson.D, db, clct string) ([]*schemas.User, error) {
	var result []*schemas.User

	// Set DB and collection
	collection := Client.Database(db).Collection(clct)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	cursor, err := collection.Find(Ctx, filter)
	cursor.Close(Ctx)

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
