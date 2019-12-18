// Provides helper functions that make interfacing with the MongoDB Go driver library easier
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

// URI for MongoDB instance
const DBuri = "mongodb://localhost:27017"

// Create Context
var Ctx context.Context

// Create MongoDB client
var Client *mongo.Client

// Error Definitions
var ErrDocumentNotFound = errors.New("No documents found under that filter.")
var ErrFindFailed = errors.New("Could not find unmarshall data. Query may be malformed")
var ErrUpdateFailed = errors.New("Could not update document.")
var ErrCursorIterationFailed = errors.New("An error occurred when iterating through a cursor")

// Wrapper for Mongo Collection
type CnctConnection struct {
	Collection *mongo.Collection
}

// Function to create new connection to Mongo Collection
func New(db, cnct string) (c CnctConnection) {
	c.Collection = Client.Database(db).Collection(cnct)
	return c
}

// Initialize connection to DB, set package wide Ctx and Client
func init() {
	log.Printf("Attempting to connect to %q", DBuri)
	conn, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clt, err := mongo.Connect(Ctx, options.Client().ApplyURI(DBuri))

	if err != nil {
		log.Panic(err)
	}
	log.Print("Connection Established")

	// Change Package level vars
	Ctx = conn
	Client = clt
}

// Simplification of collection.FindOne()
func (db CnctConnection) FindOne(filter bson.D, res interface{}) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	err = db.Collection.FindOne(Ctx, filter).Decode(res)

	// failed to unmarshall
	if err != nil {
		return ErrFindFailed
	}

	return nil
}

// Simplification of collection.Find()
func (db CnctConnection) FindMany(filter bson.D, res *[]interface{}) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	cursor, err := db.Collection.Find(Ctx, filter)

	if err != nil {
		return ErrDocumentNotFound
	}

	for cursor.Next(Ctx) {
		var doc interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			return ErrFindFailed
		}
		*res = append(*res, doc)
	}

	// unmarshall fail
	if cursor.Err() != nil {
		return ErrCursorIterationFailed
	}

	// Close cursor after we're done with it
	cursor.Close(Ctx)
	return nil
}

// Simplification of collection.UpdateOne() except it doesn't return the document
func (db CnctConnection) UpdateOne(filter, update bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	// Updated result is discarded, you can reimplement if needed
	_, err = db.Collection.UpdateOne(Ctx, filter, update)

	if err != nil {
		return ErrUpdateFailed
	}
	return nil
}

// Simplification of collection.Update() except it doesn't return a cursor
func (db CnctConnection) UpdateMany(filter, update bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	_, err = db.Collection.UpdateMany(Ctx, filter, update)

	if err != nil {
		return ErrUpdateFailed
	}
	return nil
}
