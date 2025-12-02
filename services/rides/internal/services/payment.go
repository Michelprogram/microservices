package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PaymentService struct {
	paymentServiceURL string
	client            *http.Client
}

func NewPaymentService(paymentServiceURL string) *PaymentService {
	return &PaymentService{
		paymentServiceURL: paymentServiceURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AuthorizeRequest struct {
	RideID string  `json:"ride_id"`
	Amount float64 `json:"amount"`
}

type AuthorizeResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

type CaptureRequest struct {
	PaymentID string `json:"payment_id"`
}

type CaptureResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

func (s *PaymentService) AuthorizePayment(rideID string, amount float64) (string, error) {
	url := fmt.Sprintf("%s/payments/authorize", s.paymentServiceURL)

	reqBody := AuthorizeRequest{
		RideID: rideID,
		Amount: amount,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to call payment service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(body))
	}

	var authorizeResp AuthorizeResponse
	if err := json.Unmarshal(body, &authorizeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return authorizeResp.PaymentID, nil
}

func (s *PaymentService) CapturePayment(paymentID string) error {
	url := fmt.Sprintf("%s/payments/capture", s.paymentServiceURL)

	reqBody := CaptureRequest{
		PaymentID: paymentID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call payment service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(body))
	}

	var captureResp CaptureResponse
	if err := json.Unmarshal(body, &captureResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
