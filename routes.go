package main

import (
	"net/http"

	mongo "./mongoclient"
)

func Register(w http.ResponseWriter, r *http.Request) {
	collection := mongo.client.Database("testing").Collection("numbers")
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Refresh(w http.ResponseWriter, r *http.Request) {

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

}
