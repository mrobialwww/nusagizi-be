package auth0client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ManagementClient is a thin client for Auth0 Management API.
// Used to mark email_verified=true after successful OTP verification.
// IMPORTANT: It uses M2M credentials which MUST be kept server-side and NEVER embedded in the Flutter app.
type ManagementClient struct {
	domain       string
	clientID     string
	clientSecret string
	httpClient   *http.Client

	cachedToken string
	tokenExpiry time.Time
}

func NewManagementClient(domain, clientID, clientSecret string) *ManagementClient {
	return &ManagementClient{
		domain:       domain,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// getAccessToken fetches and caches the M2M access token in memory.
// It is not strictly thread-safe; use a mutex or centralized cache for high loads.
func (c *ManagementClient) getAccessToken(ctx context.Context) (string, error) {
	if c.cachedToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.cachedToken, nil
	}

	body := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"audience":      fmt.Sprintf("https://%s/api/v2/", c.domain),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		fmt.Sprintf("https://%s/oauth/token", c.domain),
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal request token M2M: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token M2M gagal, status: %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	c.cachedToken = tr.AccessToken
	// Subtract a 60-second buffer before actual expiration.
	c.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn-60) * time.Second)

	return c.cachedToken, nil
}

type managementUser struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Identities []struct {
		Connection string `json:"connection"`
	} `json:"identities"`
}

// FindDatabaseUserByEmail finds an Auth0 user by email registered via
// the "Username-Password-Authentication" database connection.
//
// Note: This ensures we verify the main database account with a password,
// rather than the separate Passwordless "email" identity used for OTP.
func (c *ManagementClient) FindDatabaseUserByEmail(ctx context.Context, email string) (*managementUser, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		"https://%s/api/v2/users-by-email?email=%s",
		c.domain, url.QueryEscape(email),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("users-by-email gagal, status: %d", resp.StatusCode)
	}

	var users []managementUser
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	for _, u := range users {
		for _, id := range u.Identities {
			if id.Connection == "Username-Password-Authentication" {
				user := u
				return &user, nil
			}
		}
	}

	return nil, fmt.Errorf("user database dengan email %s tidak ditemukan", email)
}

// MarkEmailVerified sets email_verified=true for the specified userID.
func (c *ManagementClient) MarkEmailVerified(ctx context.Context, userID string) error {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(map[string]bool{"email_verified": true})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://%s/api/v2/users/%s", c.domain, url.PathEscape(userID))

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal update email_verified, status: %d", resp.StatusCode)
	}

	return nil
}
