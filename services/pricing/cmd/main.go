package main

import (
	"log"
	"net/http"
	"services/pricing/internal"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/price", internal.GetPriceHandler).Methods("GET")

	return r
}

func main() {
	router := SetupRouter()

	port := "4000" // port du microservice pricing
	log.Println("Pricing service running on port", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Error starting Pricing service:", err)
	}
}
