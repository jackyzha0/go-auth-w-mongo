// Defines routes for HTTP server
package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../schemas"
	db "github.com/jackyzha0/monGo-driver-wrapper"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	var creds schemas.Credentials
	decodeErr := json.NewDecoder(r.Body).Decode(&creds)

	if decodeErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, structure incorrect.\n")
		return
	}

	filter := bson.D{{"email", creds.Email}}
	var res schemas.User
	err := Users.FindOne(filter, &res)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// New Session Token
	sessionToken, _ := uuid.NewV4()

	// Update User token
	// !!! TODO

	// catch write fail
	// !!! TODO

	// Write cookie to client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(120 * time.Minute),
	})
	w.WriteHeader(http.StatusOK)

	log.Printf("User %q logged in with token %q", res.Name, sessionToken.String())
}

func Register(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
