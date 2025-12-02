package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserService struct {
	usersServiceURL string
}

func NewUserService(usersServiceURL string) *UserService {
	return &UserService{usersServiceURL: usersServiceURL}
}

func (s *UserService) GetAvailableDriver() (string, error) {
	url := fmt.Sprintf("%s/drivers?available=true", s.usersServiceURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("users service returned status %d: %s", resp.StatusCode, string(body))
	}

	var drivers []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		IsAvailable bool   `json:"is_available"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(drivers) == 0 {
		return "", fmt.Errorf("no available drivers found")
	}

	// Return the first available driver's ID
	return drivers[0].ID, nil
}

func (s *UserService) UpdateDriverStatus(driverID string, isAvailable bool) error {
	url := fmt.Sprintf("%s/drivers/%s/status", s.usersServiceURL, driverID)

	payload := struct {
		IsAvailable bool `json:"is_available"`
	}{
		IsAvailable: isAvailable,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("users service returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[UPDATE] Driver %s availability set to %v", driverID, isAvailable)
	return nil
}
