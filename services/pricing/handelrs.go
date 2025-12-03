package main

import (
	"encoding/json"
	"net/http"
)

func GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		http.Error(w, "Missing parameters 'from' and/or 'to'", http.StatusBadRequest)
		return
	}

	price, ok := GetPrice(from, to)
	if !ok {
		http.Error(w, "No price found for selected zones", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
