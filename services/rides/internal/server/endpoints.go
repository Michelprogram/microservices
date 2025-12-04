package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rides/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) createRide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PassengerID string `json:"passengerId"`
		FromZone    string `json:"from_zone"`
		ToZone      string `json:"to_zone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	price, err := s.pricingService.GetPrice(req.FromZone, req.ToZone)

	if err != nil {
		log.Printf("[WARN] Failed to get price: %v", err)
		http.Error(w, "Failed to get price", http.StatusInternalServerError)
		return
	}


	driverID, err := s.userService.GetAvailableDriver()
	if err != nil {
		log.Printf("[ERROR] Failed to get available driver: %v", err)
		http.Error(w, "No available driver found", http.StatusServiceUnavailable)
		return
	}

	rideID := primitive.NewObjectID()

	paymentID, err := s.paymentService.AuthorizePayment(rideID.Hex(), price)
	if err != nil {
		log.Printf("[WARN] Failed to authorize payment: %v", err)
		http.Error(w, "Failed to authorize payment", http.StatusInternalServerError)
	}

	ride := &types.Ride{
		ID:            rideID,
		PassengerID:   req.PassengerID,
		PaymentID:     paymentID,
		DriverID:      driverID,
		FromZone:      req.FromZone,
		ToZone:        req.ToZone,
		Price:         price,
		Status:        "ASSIGNED",
		PaymentStatus: "PENDING",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.db.CreateRide(ctx, ride)
	if err != nil {
		log.Printf("[ERROR] Failed to create ride: %v", err)
		http.Error(w, "Error creating ride", http.StatusInternalServerError)
		return
	}

	if err := s.userService.UpdateDriverStatus(driverID, false); err != nil {
		log.Printf("[WARN] Failed to update driver status: %v", err)
	}

	log.Printf("[CREATE] Nouvelle course créée: ID=%s, Passenger=%s, Driver=%s", ride.ID.Hex(), ride.PassengerID, ride.DriverID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ride)
}

func (s *Server) getRide(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ride, err := s.db.GetRideByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Ride not found", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Failed to get ride: %v", err)
		http.Error(w, "Error retrieving ride", http.StatusInternalServerError)
		return
	}

	log.Printf("[READ] Course récupérée: ID=%s", idStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ride)
}

func (s *Server) updateRideStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.db.UpdateRideStatus(ctx, id, req.Status)
	if err != nil {
		log.Printf("[ERROR] Failed to update ride status: %v", err)
		http.Error(w, "Error updating ride status", http.StatusInternalServerError)
		return
	}

	// If status is COMPLETED, capture payment
	if req.Status == "COMPLETED" {
		ride, err := s.db.GetRideByID(ctx, id)
		if err == nil && ride.PaymentID != "" {
			err = s.paymentService.CapturePayment(ride.PaymentID)
			if err != nil {
				log.Printf("[ERROR] Failed to capture payment: %v", err)
			} else {
				err = s.db.UpdateRidePaymentStatus(ctx, id, "CAPTURED")
				if err != nil {
					log.Printf("[ERROR] Failed to update payment status: %v", err)
				}
			}
		}

		if err := s.userService.UpdateDriverStatus(ride.DriverID, true); err != nil {
			log.Printf("[WARN] Failed to update driver status: %v", err)
		}
	}

	// Get updated ride to return
	ride, err := s.db.GetRideByID(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Failed to get updated ride: %v", err)
		http.Error(w, "Error retrieving updated ride", http.StatusInternalServerError)
		return
	}

	log.Printf("[UPDATE] Statut de la course %s mis à jour: %s", idStr, req.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ride)
}
