package db

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_uri = "mongodb://localhost:27017"

// Create DB Connection
var Ctx context.Context
var Client *mongo.Client

var ErrDocumentNotFound = errors.New("no found under that filter")
var ErrCouldNotUnMarshall = errors.New("could not find unmarshall data. query may be malformed")
var ErrCursorIterationFailed = errors.New("an error occurred when iterating through a cursor")
var ErrCollectionNotDefined = errors.New("collection not defined in db.go")

type Collection struct {
	DB         string
	Collection string
}

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

func (db Collection) FindOne(filter bson.D, res interface{}) (err error) {
	// Set DB and collection
	collection := Client.Database(db.DB).Collection(db.Collection)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Find Operation
	err = collection.FindOne(Ctx, filter).Decode(res)

	// failed to unmarshall
	if err != nil {
		return ErrCouldNotUnMarshall
	}

	return nil
}

func (db Collection) FindMany(filter bson.D, res *[]interface{}) (err error) {
	// Set DB and collection
	collection := Client.Database(db.DB).Collection(db.Collection)
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Find Operation
	cursor, err := collection.Find(Ctx, filter)

	if err != nil {
		return ErrDocumentNotFound
	}

	for cursor.Next(Ctx) {
		var doc interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			return ErrCouldNotUnMarshall
		}
		*res = append(*res, doc)
	}

	// unmarshall fail
	if cursor.Err() != nil {
		log.Fatalf("cursor iteration failed with filter %q", filter)
		return ErrCursorIterationFailed
	}

	// Close cursor after we're done with it
	cursor.Close(Ctx)
	return nil
}
