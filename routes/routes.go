package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../db"

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

func Register(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
	var user schemas.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request, structure incorrect.\n")
		return
	}

	filter := bson.M{"email": user.Email}
}

func Refresh(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
