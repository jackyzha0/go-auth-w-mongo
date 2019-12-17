package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"./db"

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
	// collection := db.Client.Database("signin").Collection("students")
	// db.Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	// cur, err := collection.find(db.Ctx, bson.D{"email":""})
	//
	// if err != nil {
	//
	// }
	//
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Hello!\n")
}

func Refresh(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
