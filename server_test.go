package main

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	m "go.mongodb.org/mongo-driver/mongo"

	"./db"
)

// Create new connection to Users Collection. To ensure test state, ensure that a Database called `exampleDB` and a Collection called `test` exists and contains two documents
// {name: "bob", surname: "joe"}
// {name: "sally", surname: "joe"}
var TestCollection = db.New("exampleDB", "test")

// Check to see if JSON content is the same
func assertJSON(t *testing.T, got interface{}, want string) {
	t.Helper()
	got = fmt.Sprintf("%v", got)
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Driver to test Find
func TestFind(t *testing.T) {
	t.Run("Find single element present", func(t *testing.T) {
		// should return single document even when multiple satisfy query
		filter := bson.D{{"name", "bob"}}
		var res interface{}
		err := TestCollection.FindOne(filter, &res)

		want := "[{_id ObjectID(\"5dfa95498e08ca409f93d998\")} {name bob} {surname joe}]"

		assertNoError(t, err)
		assertJSON(t, res, want)
	})

	t.Run("Find single not present", func(t *testing.T) {
		// should error
		filter := bson.D{{"name", "john"}}
		var res interface{}
		err := TestCollection.FindOne(filter, &res)
		assertError(t, err, m.ErrNoDocuments)
	})

	t.Run("Find with invalid field", func(t *testing.T) {
		// should error
		filter := bson.D{{"birthday", "04/16/01"}}
		var res interface{}
		err := TestCollection.FindOne(filter, &res)
		assertError(t, err, m.ErrNoDocuments)
	})

	t.Run("Find multiple present", func(t *testing.T) {
		// should return all documents satisfying query
		filter := bson.D{{"surname", "joe"}}
		var res []interface{}
		err := TestCollection.FindMany(filter, &res)

		want := "[[{_id ObjectID(\"5dfa95498e08ca409f93d998\")} {name bob} {surname joe}] [{_id ObjectID(\"5dfa95598e08ca409f93d999\")} {name sally} {surname joe}]]"

		assertNoError(t, err)
		assertJSON(t, res, want)
	})

	t.Run("Find multiple not present", func(t *testing.T) {
		// should return an empty cursor
		filter := bson.D{{"name", "john"}}
		var res []interface{}
		err := TestCollection.FindMany(filter, &res)

		assertNoError(t, err)
		if len(res) != 0 {
			t.Fatalf("got %+v, wanted []", res)
		}
	})
}

// // Driver to test InsertOne
// func TestInsertOne(t *testing.T) {
//
// }
//
// // Driver to test InsertMany
// func TestInsertMany(t *testing.T) {
// }
//
// // Driver to test DeleteOne
// func TestDeleteOne(t *testing.T) {
// }
//
// // Driver to test DeleteMany
// func TestDeleteMany(t *testing.T) {
// }
//
// // Driver to test UpdateOne
// func TestUpdateOne(t *testing.T) {
// }
//
// // Driver to test UpdateMany
// func TestUpdateMany(t *testing.T) {
// }
