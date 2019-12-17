package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	mux "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Create DB Connection
var ctx context.Context
var client *mongo.Client

func main() {
	// Initialize DB
	initDB()

	// Define Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/register", Register)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/refresh", Refresh)
	r.HandleFunc("/dashboard", Dashboard)
	r.HandleFunc("/dbhealthcheck", HealthCheck)
	http.Handle("/", r)

	// Start HTTP server
	server := newServer(":"+strconv.Itoa(8080), r)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func initDB() {
	conn, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clt, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Panic(err)
	}

	// Change Package level vars
	ctx = conn
	client = clt
}

func newServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := client.Ping(ctx, readpref.Primary())

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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello!\n")
}

func Refresh(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
