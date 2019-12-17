package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../db"
	"../schemas"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

func Login(w http.ResponseWriter, r *http.Request) {
	var creds schemas.Credentials
	decodeErr := json.NewDecoder(r.Body).Decode(&creds)

	if decodeErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, structure incorrect.\n")
		return
	}

	filter := bson.M{"email": creds.Email}
	var res schemas.User
	database := db.Collection{DB: "exampleDB", Collection: "users"}
	err := database.FindOne(filter, &res)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// New Session Token
	sessionToken, _ := uuid.NewV4()

	// Update User token

	// catch write fail
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error occured when attempting to write session token to database")
		return
	}

	w.WriteHeader(http.StatusOK)
	// Write cookie to client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(120 * time.Second),
	})

	log.Printf("User %q logged in with token %q", res.Name, sessionToken.String())
}

func Register(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
