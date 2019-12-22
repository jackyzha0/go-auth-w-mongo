// Defines routes for the app
package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackyzha0/go-auth-w-mongo/schemas"
	db "github.com/jackyzha0/monGo-driver-wrapper"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

// token status
const (
	TokenValid int = 0
	TokenBadStruct int = 1
	TokenExpired int = 2
	TokenNotAdmin int = 3
	TokenOtherErr int = -1
)

// Users is a new connection to Users Collection
var Users = db.New("mongodb://localhost:27017", "exampleDB", "users")

// hashes and salts current password
func hashSalt(pass []byte) (string, error) {
    hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

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

// check if valid token
func isValidToken(r *http.Request, adminCheck bool) (status int, retUser schemas.User) {
	c, cookieFetchErr := r.Cookie("session_token")

	// not auth'ed, redirect to login
	if cookieFetchErr != nil {
		if cookieFetchErr == http.ErrNoCookie {
			return TokenExpired, schemas.User{}
		}
		return TokenBadStruct, schemas.User{}
	}

	// if no err, get cookie value
	sessionToken := c.Value

	filter := bson.D{{"sessionToken", sessionToken}}
	var res schemas.User
	findErr := Users.FindOne(filter, &res)

	if findErr != nil {

		// no user with matching session_token
		if findErr == mongo.ErrNoDocuments {
			return TokenExpired, schemas.User{}
		}

		// other error
		return TokenOtherErr, schemas.User{}
	}

	// parse time
	expireTime, timeParseErr := time.Parse(time.RFC3339, res.SessionExpires)

	// token time invalid
	if timeParseErr != nil {
		return TokenBadStruct, schemas.User{}
	}

	// token expired
	if time.Now().After(expireTime) {
		return TokenExpired, schemas.User{}
	}

	if adminCheck && !res.IsAdmin {
		return TokenNotAdmin, schemas.User{}
	}

	// token ok
	return TokenValid, res
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
	return
}

// Register is the endpoint to register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	stat, _ := isValidToken(r, true)

	// check if error or not authed
	switch stat {
	case TokenExpired:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	case TokenBadStruct:
		w.WriteHeader(http.StatusBadRequest)
		return
	case TokenNotAdmin:
		w.WriteHeader(http.StatusUnauthorized)
		return
	case TokenOtherErr:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := new(schemas.User)

	// parse Form request
	parseErr := r.ParseForm()
	decoder := schema.NewDecoder()
	parseErr = decoder.Decode(newUser, r.PostForm)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, structure incorrect.\n")
		return
	}

	// create non-admin
	newUser.IsAdmin = false

	// hash password
    hash, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
    if hashErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
    }
    newUser.Password = string(hash)

	// insert document
	ID, insertErr := Users.InsertOne(newUser)

	// email already exists
	if insertErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, user with that email exists.\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Created new user with ID %s.\n", ID)
	return
}

// Dashboard is the endpoint to display a welcome page to auth'd users
func Dashboard(w http.ResponseWriter, r *http.Request) {
	stat, res := isValidToken(r, false)

	switch stat {
	case TokenExpired:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	case TokenBadStruct:
		w.WriteHeader(http.StatusBadRequest)
		return
	case TokenOtherErr:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// user is authed! display data or do something
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome back %s!\n", res.Name)
	return
}
