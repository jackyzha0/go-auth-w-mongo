package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackyzha0/go-auth-w-mongo/schemas"
	db "github.com/jackyzha0/monGo-driver-wrapper"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gorilla/schema"
)

// Create new connection to Users Collection
var Users = db.New("mongodb://localhost:27017", "exampleDB", "users")

// Endpoint to check if connection to database is healthy
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db.Ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err := db.Client.Ping(db.Ctx, readpref.Primary())

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Waited 2 seconds, could not connect to server.\n")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Connection to Database established.\n")
	}
}

// Endpoint to users to login through, injects a sessionToken that is valid for 2 hours
func Login(w http.ResponseWriter, r *http.Request) {
	creds := new(schemas.Credentials)

	// parse Form request
	parseErr := r.ParseForm()
	decoder := schema.NewDecoder()
	// r.PostForm is a map of our POST form values
	parseErr = decoder.Decode(creds, r.PostForm)

	if parseErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, structure incorrect. Please include a password and email.\n")
		return
	}

	filter := bson.D{{"email", creds.Email}}
	var res schemas.User
	findErr := Users.FindOne(filter, &res)

	if findErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// New Session Token
	sessionToken := uuid.NewV4()
	expiry := time.Now().Add(120 * time.Minute)
	expiryStr := expiry.Format(time.RFC3339)
	// to parse time back, do
	// t, err := time.Parse(time.RFC3339, str)

	// Update User
	update := bson.D{{
		"$set", bson.D{
			{"session_token", sessionToken.String()},
			{"session_expires", expiryStr}}}}
	_,_,updateErr := Users.UpdateOne(filter, update)

	if updateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write cookie to client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiry,
	})
	w.WriteHeader(http.StatusOK)

	log.Printf("User %q logged in with token %q", res.Name, sessionToken.String())
}

// Endpoint to register a new user
func Register(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
