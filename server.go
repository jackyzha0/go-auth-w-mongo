package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"./routes"

	mux "github.com/gorilla/mux"
)

const port = 8080

func main() {
	// Define Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/register", routes.Register)
	r.HandleFunc("/login", routes.Login)
	r.HandleFunc("/dashboard", routes.Dashboard)
	r.HandleFunc("/qtest", routes.Test)
	r.HandleFunc("/dbhealthcheck", routes.HealthCheck)
	http.Handle("/", r)

	// Start HTTP server
	server := newServer(":"+strconv.Itoa(port), r)
	log.Printf("Starting server on %d", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func newServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
}
