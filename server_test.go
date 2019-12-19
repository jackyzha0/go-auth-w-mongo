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
	var val []Doc
	for _, el := range got {
		var d Doc
		bsonBytes, _ := bson.Marshal(el)
		bson.Unmarshal(bsonBytes, &d)
		val = append(val, d)
	}

	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %+v want %+v", got, want)
	}
}

func assertMatchMod(t *testing.T, gotmatch, gotmod int64, wantmatch, wantmod int) {
	if !(gotmatch == int64(wantmatch) && gotmod == int64(wantmod)) {
		t.Errorf("wanted %d matches and %d modification, got %d matches and %d modifications", wantmatch, wantmod, gotmatch, gotmod)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got an error %q but didn't want one", got)
	}
}

func assertZeroLength(t *testing.T, got []interface{}) {
	t.Helper()
	if len(got) != 0 {
		t.Error("Expected result to be empty, was not")
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	if got == nil {
		t.Errorf("didn't get an error but wanted %q", want)
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
		assertZeroLength(t, res)
	})
}

func TestInsert(t *testing.T) {
	t.Run("Insert single valid document", func(t *testing.T) {
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

	t.Run("Insert many valid documents", func(t *testing.T) {
		john := Doc{"alex", "smith"}
		betty := Doc{"alex", "hansen"}
		sl := []interface{}{john, betty}
		err := TestCollection.InsertMany(sl)
		assertNoError(t, err)

		filter := bson.D{{"name", "alex"}}
		var res []interface{}
		err = TestCollection.FindMany(filter, &res)

		want := []Doc{Doc{"alex", "smith"}, Doc{"alex", "hansen"}}

		assertNoError(t, err)
		assertMultipleDoc(t, res, want)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update single doc set existing field", func(t *testing.T) {
		rebecca := Doc{"rebecca", "joe"}
		err := TestCollection.InsertOne(rebecca)
		assertNoError(t, err)

		filter := bson.D{{"name", "rebecca"}, {"surname", "o'connor"}}
		var res Doc
		err = TestCollection.FindOne(filter, &res)
		assertError(t, err, m.ErrNoDocuments)

		update_filter := bson.D{{"name", "rebecca"}}
		update := bson.D{{"$set", bson.D{{"surname", "o'connor"}}}}
		err, match, mod := TestCollection.UpdateOne(update_filter, update)
		assertMatchMod(t, match, mod, 1, 1)
		assertNoError(t, err)

		err = TestCollection.FindOne(filter, &res)
		assertNoError(t, err)
	})

	t.Run("update single doc set nonexistant field", func(t *testing.T) {
		kevin := Doc{"kevin", "kwon"}
		err := TestCollection.InsertOne(kevin)
		assertNoError(t, err)

		filter := bson.D{{"age", 24}}
		var res Doc
		err = TestCollection.FindOne(filter, &res)
		assertError(t, err, m.ErrNoDocuments)

		update_filter := bson.D{{"name", "kevin"}, {"surname", "kwon"}}
		update := bson.D{{"$set", bson.D{{"age", 24}}}}
		err, match, mod := TestCollection.UpdateOne(update_filter, update)
		assertMatchMod(t, match, mod, 1, 1)
		assertNoError(t, err)

		err = TestCollection.FindOne(filter, &res)
		assertNoError(t, err)
	})

	t.Run("update single doc multiple matches", func(t *testing.T) {
		jessica := Doc{"jessica", "wu"}
		err := TestCollection.InsertOne(jessica)
		assertNoError(t, err)

		update_filter := bson.D{{"surname", "joe"}}
		update := bson.D{{"$set", bson.D{{"surname", "o'connor"}}}}
		err, match, mod := TestCollection.UpdateOne(update_filter, update)
		assertMatchMod(t, match, mod, 2, 1)
		assertNoError(t, err)
	})

	t.Run("update multiples docs set existing field", func(t *testing.T) {

	})

	t.Run("update multiple docs set nonexistant field", func(t *testing.T) {

	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete single valid document", func(t *testing.T) {
		albert := Doc{"albert", "yip"}
		_ = TestCollection.InsertOne(albert)
		filter := bson.D{{"name", "albert"}, {"surname", "yip"}}

		var res interface{}
		err := TestCollection.FindOne(filter, &res)
		assertNoError(t, err)

		err = TestCollection.DeleteOne(filter)
		assertNoError(t, err)

		var res1 interface{}
		err = TestCollection.FindOne(filter, &res1)
		assertError(t, err, m.ErrNoDocuments)
	})

	t.Run("Delete many valid documents", func(t *testing.T) {
		nick := Doc{"nick", "zheng"}
		stephen := Doc{"stephen", "zheng"}
		sl := []interface{}{nick, stephen}
		err := TestCollection.InsertMany(sl)
		assertNoError(t, err)

		filter := bson.D{{"surname", "zheng"}}
		var res []interface{}
		err = TestCollection.FindMany(filter, &res)
		want := []Doc{Doc{"nick", "zheng"}, Doc{"stephen", "zheng"}}
		assertNoError(t, err)
		assertMultipleDoc(t, res, want)

		err = TestCollection.DeleteMany(filter)
		assertNoError(t, err)

		var nres []interface{}
		err = TestCollection.FindMany(filter, &nres)
		assertNoError(t, err)
		assertZeroLength(t, nres)
	})
}
