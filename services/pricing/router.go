package main

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/price", GetPriceHandler).Methods("GET")

	return r
}
