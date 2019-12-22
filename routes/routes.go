package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/jackyzha0/go-auth-w-mongo/schemas"
	db "github.com/jackyzha0/monGo-driver-wrapper"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gorilla/schema"
)

// Users is a new connection to Users Collection
var Users = db.New("mongodb://localhost:27017", "exampleDB", "users")

// refresh/set user token by email
func refreshToken(email string) (c *http.Cookie, ok bool) {
	// New Session Token
	sessionToken := uuid.NewV4()
	expiry := time.Now().Add(120 * time.Minute)
	expiryStr := expiry.Format(time.RFC3339)

	// Update User
	filter := bson.D{{"email", email}}
	update := bson.D{{
		"$set", bson.D{
			{"sessionToken", sessionToken.String()},
			{"sessionExpires", expiryStr}}}}
	_,_,updateErr := Users.UpdateOne(filter, update)

	if updateErr != nil {
		return nil, false
	}

	return &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiry,
	}, true
}


// HealthCheck is an endpoint to check if connection to database is healthy
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

// Login is an endpoint to users to login through, injects a sessionToken that is valid for 2 hours
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

	c, ok := refreshToken(res.Email)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write cookie to client
	http.SetCookie(w, c)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

// Register is the endpoint to register a new user
func Register(w http.ResponseWriter, r *http.Request) {
}

// Dashboard is the endpoint to display a welcome page to auth'd users
func Dashboard(w http.ResponseWriter, r *http.Request) {
	c, cookieFetchErr := r.Cookie("session_token")

	// not auth'ed, redirect to login
	if cookieFetchErr != nil {
		if cookieFetchErr == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if no err, get cookie value
	sessionToken := c.Value

	filter := bson.D{{"sessionToken", sessionToken}}
	log.Printf("filter -> %+v", filter)
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

	log.Printf("parsed struct %+v", res)
	expireTime, timeParseErr := time.Parse(time.RFC3339, res.SessionExpires)

	if timeParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// token expired
	if time.Now().After(expireTime) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// user is authed! 
}
