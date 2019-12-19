package main

import (
	"os"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	m "go.mongodb.org/mongo-driver/mongo"

	"./db"
)

// Create new connection to Users Collection.
// {name: "bob", surname: "joe"}
// {name: "sally", surname: "joe"}
var TestCollection = db.New("exampleDB", "test")

type Doc struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// Check to see if JSON content is the same
func assertSingleDoc(t *testing.T, got Doc, want Doc) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertMultipleDoc(t *testing.T, got []interface{}, want []Doc) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
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

func TestMain(m *testing.M) {
	// setup, add two docs
	bob := Doc{"bob", "joe"}
	sally := Doc{"sally", "joe"}
	sl := []interface{}{bob, sally}
	TestCollection.InsertMany(sl)

	// run rest of tests
	exitVal := m.Run()

	// cleanup, drop test collection
	TestCollection.Drop()

	os.Exit(exitVal)
}

// Driver to test Find
func TestFind(t *testing.T) {
	t.Run("Find single element present", func(t *testing.T) {
		// should return single document even when multiple satisfy query
		filter := bson.D{{"name", "bob"}}
		var res Doc
		err := TestCollection.FindOne(filter, &res)

		want := Doc{"bob", "joe"}

		assertNoError(t, err)
		assertSingleDoc(t, res, want)
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

		want := []Doc{Doc{"bob", "joe"}, Doc{"sally", "joe"}}

		assertNoError(t, err)
		assertMultipleDoc(t, res, want)
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

func TestInsert(t *testing.T) {
	t.Run("insert single valid document", func(t *testing.T) {
		// make sure john doesnt exist before
		filter := bson.D{{"name", "john"}}
		var res interface{}
		err := TestCollection.FindOne(filter, &res)
		assertError(t, err, m.ErrNoDocuments)

		john := Doc{"john", "smith"}
		err = TestCollection.InsertOne(john)
		assertNoError(t, err)

		// make sure john exists now
		filter = bson.D{{"name", "john"}}
		err = TestCollection.FindOne(filter, &res)
		assertNoError(t, err)
	})

	// t.Run("insert many valid documents", func(t *testing.T) {
	//   john := Doc{"john", "smith"}
	//   betty := Doc{"betty", "hansen"}
	//   sl := []Doc{john, betty}
	//   err = TestCollection.InsertMany(sl)
	//   assertNoError(t, err)
	// })
}
