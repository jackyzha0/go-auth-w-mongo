// Provides helper functions that make interfacing with the MongoDB Go driver library easier
package db

import (
	"context"
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

const OperationTimeOut = 5

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
	conn, _ := context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
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
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	err = db.Collection.FindOne(Ctx, filter).Decode(res)

	// failed to unmarshall
	if err != nil {
		return err
	}

	return nil
}

// Simplification of collection.Find()
func (db CnctConnection) FindMany(filter bson.D, res *[]interface{}) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	cursor, err := db.Collection.Find(Ctx, filter)

	if err != nil {
		return err
	}

	for cursor.Next(Ctx) {
		var doc interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			return err
		}
		*res = append(*res, doc)
	}

	// unmarshall fail
	if cursor.Err() != nil {
		return err
	}

	// Close cursor after we're done with it
	cursor.Close(Ctx)
	return nil
}

// Simplification of collection.UpdateOne() except it doesn't return the document
func (db CnctConnection) UpdateOne(filter, update bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	// Updated result is discarded, you can reimplement if needed
	_, err = db.Collection.UpdateOne(Ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}

// Simplification of collection.Update() except it doesn't return a cursor
func (db CnctConnection) UpdateMany(filter, update bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	_, err = db.Collection.UpdateMany(Ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}

// Simplification of InsertOne, doesn't return document and accepts arbitrary structs
func (db CnctConnection) InsertOne(str interface{}) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	_, err = db.Collection.InsertOne(Ctx, str)
	if err != nil {
		return err
	}
	return nil
}

// Simplification of InsertMany, takes slice of structs
func (db CnctConnection) InsertMany(str []interface{}) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	_, err = db.Collection.InsertMany(Ctx, str)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper for collection.DeleteOne()
func (db CnctConnection) DeleteOne(filter bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err = db.Collection.DeleteOne(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper for collection.DeleteMany()
func (db CnctConnection) DeleteMany(filter bson.D) (err error) {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err = db.Collection.DeleteMany(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
