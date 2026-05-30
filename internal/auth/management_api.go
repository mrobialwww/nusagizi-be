package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// tokenCache provides thread-safe in-memory storage for the Auth0 M2M token.
// This prevents redundant token requests, reducing latency and avoiding rate limits.
type tokenCache struct {
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

var cache = &tokenCache{}

// getM2MToken retrieves a Machine-to-Machine (M2M) access token via the Client Credentials Grant.
// It returns the cached token if valid, or fetches a new one from Auth0 if expired.
func getM2MToken(domain, clientID, clientSecret string) (string, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Return cached token if it has at least 30 seconds of validity remaining.
	if cache.token != "" && time.Now().Before(cache.expiresAt.Add(-30*time.Second)) {
		return cache.token, nil
	}

	// Prepare the token request payload for the Auth0 /oauth/token endpoint.
	audience := fmt.Sprintf("https://%s/api/v2/", domain)
	body, _ := json.Marshal(map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
	})

	// Execute the HTTP POST request to obtain the token.
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to request M2M token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("M2M token request failed with status: %d", resp.StatusCode)
	}

	// Decode the JSON response to extract the token and its lifespan.
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"` // Token lifespan in seconds
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode M2M token response: %w", err)
	}

	// Persist the new token and calculate its exact expiration timestamp.
	cache.token = result.AccessToken
	cache.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return cache.token, nil
}

// AssignRoleToUser binds a specific Auth0 Role to a User via the Management API.
// It abstracts the M2M authentication process and executes the role assignment payload.
func AssignRoleToUser(domain, clientID, clientSecret, auth0Sub, roleID string) error {
	// Ensure an active M2M token is available for authorization.
	token, err := getM2MToken(domain, clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("could not get M2M token: %w", err)
	}

	// Construct the JSON payload containing the array of roles to assign.
	body, _ := json.Marshal(map[string][]string{
		"roles": {roleID},
	})

	// Initialize the POST request to the Auth0 Management API users endpoint.
	url := fmt.Sprintf("https://%s/api/v2/users/%s/roles", domain, auth0Sub)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to build assign-role request: %w", err)
	}
	
	// Inject required headers, including the Bearer token for authentication.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Dispatch the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call assign-role API: %w", err)
	}
	defer resp.Body.Close()

	// Auth0 returns HTTP 204 (No Content) upon successful role assignment.
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("assign-role API returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}
