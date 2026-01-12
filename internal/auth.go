package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
)

type OpenIDConfiguration struct {
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
	TokenEndpoint               string `json:"token_endpoint"`
}

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func GetOpenIDConfiguration(issuer string) (*OpenIDConfiguration, error) {
	// Ensure issuer doesn't have a trailing slash for consistency
	issuer = strings.TrimSuffix(issuer, "/")
	configURL := fmt.Sprintf("%s/.well-known/openid-configuration", issuer)

	resp, err := http.Get(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch openid configuration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch openid configuration: status %d", resp.StatusCode)
	}

	var config OpenIDConfiguration
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode openid configuration: %w", err)
	}

	if config.DeviceAuthorizationEndpoint == "" {
		return nil, fmt.Errorf("oidc provider does not support device authorization grant")
	}

	return &config, nil
}

func DeviceAuthorize(endpoint, clientID, audience, scope string) (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", scope)
	data.Set("audience", audience)

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("device authorization failed: %s", string(body))
	}

	var deviceCode DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return nil, err
	}

	return &deviceCode, nil
}

func PollForToken(ctx context.Context, endpoint, clientID, deviceCode string, interval int) (*TokenResponse, error) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientID)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			resp, err := http.Post(
				endpoint,
				"application/x-www-form-urlencoded",
				strings.NewReader(data.Encode()),
			)
			if err != nil {
				return nil, err
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var token TokenResponse
				if err := json.Unmarshal(body, &token); err != nil {
					return nil, err
				}
				return &token, nil
			}

			var tokenErr TokenError
			if err := json.Unmarshal(body, &tokenErr); err != nil {
				return nil, fmt.Errorf("unexpected response: %s", string(body))
			}

			if tokenErr.Error == "authorization_pending" {
				continue
			}

			if tokenErr.Error == "slow_down" {
				ticker.Reset(time.Duration(interval+5) * time.Second)
				continue
			}

			return nil, fmt.Errorf("token error: %s - %s", tokenErr.Error, tokenErr.ErrorDescription)
		}
	}
}

func LoginWithDeviceFlow(issuer, clientID, audience, scope string) (*TokenResponse, error) {
	color.Cyan("Discovering OIDC configuration from %s...", issuer)
	oidcConfig, err := GetOpenIDConfiguration(issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC configuration: %w", err)
	}

	deviceCode, err := DeviceAuthorize(oidcConfig.DeviceAuthorizationEndpoint, clientID, audience, scope)
	if err != nil {
		return nil, fmt.Errorf("device authorization failed: %w", err)
	}

	color.Yellow("\nTo authenticate, visit:")
	color.Cyan("  %s", deviceCode.VerificationURIComplete)
	color.Yellow("\nOr enter this code at %s:", deviceCode.VerificationURI)
	color.Green("  %s", deviceCode.UserCode)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deviceCode.ExpiresIn)*time.Second)
	defer cancel()

	color.Yellow("Waiting for authentication...")
	token, err := PollForToken(ctx, oidcConfig.TokenEndpoint, clientID, deviceCode.DeviceCode, deviceCode.Interval)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	color.Green("✓ Authentication successful!")
	return token, nil
}
