package main

import (
	"net/http"
	"strconv"
	"time"

	mux "github.com/gorilla/mux"
)

func main() {
	// Define Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/register", Register)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/refresh", Refresh)
	r.HandleFunc("/dashboard", Dashboard)
	http.Handle("/", r)

	// Start HTTP server
	server := newServer(":"+strconv.Itoa(8080), r)
	panic(server.ListenAndServe())
}

func newServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
}
