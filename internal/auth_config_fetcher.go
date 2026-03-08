package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthConfigResponse struct {
	Issuer   string `json:"issuer"`
	ClientID string `json:"client_id"`
	Audience string `json:"audience"`
}

func FetchAuthConfig(serverURL string) (*AuthConfigResponse, error) {
	url := fmt.Sprintf("%s/v1/client/config/auth", serverURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch auth config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch auth config: status %d", resp.StatusCode)
	}

	var config AuthConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode auth config: %w", err)
	}

	return &config, nil
}
