package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jackyzha0/go-auth-w-mongo/middleware"
	"github.com/jackyzha0/go-auth-w-mongo/routes"

	mux "github.com/gorilla/mux"
)

// Port to run application
const port = 8080

func cleanup() {
	log.Print("Shutting down server...")
}

// Define router and start server
func main() {
	// Sigint trapper
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	// Define Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/register", middleware.Auth(routes.Register, true))
	r.HandleFunc("/login", routes.Login)
	r.HandleFunc("/dashboard", middleware.Auth(routes.Dashboard, false))
	r.HandleFunc("/logout", middleware.Auth(routes.Logout, false))

	http.Handle("/", r)

	// Start HTTP server
	server := newServer(":"+strconv.Itoa(port), r)
	log.Printf("Starting server on %d", port)
	defer cleanup()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Function to create new HTTP server
func newServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
}
