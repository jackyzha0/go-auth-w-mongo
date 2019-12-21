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
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/schema"
)

// Create new connection to Users Collection
var Users = db.New("mongodb://localhost:27017", "exampleDB", "users")

// Endpoint to check if connection to database is healthy
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	db.Ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err := db.Client.Ping(db.Ctx, mongo.readpref.Primary())

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
	log.Printf("User %q logged in with token %q", res.Name, sessionToken.String())
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

// Endpoint to register a new user
func Register(w http.ResponseWriter, r *http.Request) {
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	c, cookieFetchErr := r.Cookie("session_token")

	// not auth'ed, redirect to login
	if cookieFetchErr != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if no err, get cookie value
	sessionToken := c.Value

	filter := bson.D{{"session_token", sessionToken}}
	var res schemas.User
	findErr := Users.FindOne(filter, &res)

	if findErr != nil {

		// no user with matching session_token
		if findErr == mongo.ErrNoDocuments {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// other error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expireTime, timeParseErr := time.Parse(time.RFC3339, res.SessionExpires)

	// token expired
	if time.Now().After(expireTime) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// user is authed! 
}
