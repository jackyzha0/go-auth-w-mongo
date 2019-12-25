// Package routes defines routes for the app
package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackyzha0/go-auth-w-mongo/schemas"
	uuid "github.com/satori/go.uuid"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

var session, _ = mgo.Dial("mongodb://localhost:27017")

// Users is a new connection to Users Collection
var Users = session.DB("exampleDB").C("Users")

// refresh/set user token by email
func refreshToken(email string) (c *http.Cookie, ok bool) {
	// New Session Token
	sessionToken := uuid.NewV4()
	expiry := time.Now().Add(120 * time.Minute)
	expiryStr := expiry.Format(time.RFC3339)

	// Update User
	update := bson.M{
		"$set": bson.M{"sessionToken": sessionToken.String(),
			"sessionExpires": expiryStr}}
	updateErr := Users.Update(bson.M{"email": email}, update)

	if updateErr != nil {
		return nil, false
	}

	return &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiry,
	}, true
}

// Login is an endpoint to users to login through, injects a sessionToken that is valid for 2 hours
func Login(w http.ResponseWriter, r *http.Request) {
	creds := new(schemas.Credentials)

	parseFormErr := r.ParseForm()
	if parseFormErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", parseFormErr)
		return
	}

	decoder := schema.NewDecoder()
	parseErr := decoder.Decode(creds, r.PostForm)
	if parseErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", parseErr)
		return
	}

	filter := bson.M{"email": creds.Email}
	var res schemas.User
	findErr := Users.Find(filter).One(&res)

	// can't find user, email doesnt exist in db
	if findErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check entered password with db password
	gotPass := []byte(creds.Password)
	dbPass := []byte(res.Password)
	compErr := bcrypt.CompareHashAndPassword(dbPass, gotPass)

	// if mismatched, return unauth status
	if compErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// auth ok, refresh/set user token
	c, ok := refreshToken(res.Email)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write cookie to client
	http.SetCookie(w, c)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Register is the endpoint to register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	newUser := new(schemas.User)
	parseFormErr := r.ParseForm()
	if parseFormErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", parseFormErr)
		return
	}

	decoder := schema.NewDecoder()
	parseErr := decoder.Decode(newUser, r.PostForm)
	if parseErr != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", parseErr)
		return
	}

	// hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
	if hashErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newUser.Password = string(hash)

	// insert document
	insertErr := Users.Insert(newUser)

	// email already exists
	if insertErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, user with that email exists.\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Created new user")
}

// Dashboard is the endpoint to display a welcome page to auth'd users
func Dashboard(w http.ResponseWriter, r *http.Request) {
	e := r.Header.Get("X-res-email")
	var res schemas.User
	_ = Users.Find(bson.M{"email": e}).One(&res)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome back %s!\n", res.Name)
}
