package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type priceRequest struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Price float64 `json:"price"`
}

type PricingService struct {
	pricingServiceURL string
	client            *http.Client
}

func NewPricingService(pricingServiceURL string) *PricingService {
	return &PricingService{
		pricingServiceURL: pricingServiceURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ps *PricingService) GetPrice(fromZone, toZone string) (float64, error) {

	url := fmt.Sprintf("%s/price?from=%s&to=%s", ps.pricingServiceURL, fromZone, toZone)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call pricing service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("pricing service returned status %d: %s", resp.StatusCode, string(body))
	}

	var res priceRequest

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.Price, nil

}
