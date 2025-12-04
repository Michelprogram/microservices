package internal

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
)

type PriceResponse struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Price float64 `json:"price"`
}

func GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		http.Error(w, "Missing parameters 'from' and/or 'to'", http.StatusBadRequest)
		return
	}
	log.Println("Getting price for zones:", from, "and", to)

	response := PriceResponse{
		From:  from,
		To:    to,
		Price: randomPrice(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Récupère le prix des zones
func randomPrice() float64 {
	return rand.Float64() * 100
}
