package main

import (
	"log"
	"net/http"
)

func main() {
	router := SetupRouter()

	port := "4000" // port du microservice pricing
	log.Println("Pricing service running on port", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Error starting Pricing service:", err)
	}
}
