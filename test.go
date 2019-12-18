package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"

	. "./db"
)

// Create new connection to Users Collection
var TestCollection = New("exampleDB", "test")

// Driver to test FindOne
func FindOneTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test FindMany
func FindManyTest(w http.ResponseWriter, r *http.Request) {
	filter := bson.D{{}}
	var res []interface{}
	err := TestCollection.FindMany(filter, &res)

	if err != nil {
		log.Panic(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Results found %+v.\n", res)
}

// Driver to test UpdateOne
func UpdateOneTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test UpdateMany
func UpdateManyTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test DeleteOne
func DeleteOneTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test DeleteMany
func DeleteManyTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test InsertOne
func InsertOneTest(w http.ResponseWriter, r *http.Request) {
}

// Driver to test InsertMany
func InsertManyTest(w http.ResponseWriter, r *http.Request) {
}
