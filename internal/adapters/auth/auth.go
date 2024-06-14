package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"hex/internal/application/auth"
)

type railsAuthService struct {
	railsBaseURL string
}

func NewRailsAuthService(railsBaseURL string) auth.AuthService {
	return &railsAuthService{railsBaseURL: railsBaseURL}
}

func (s *railsAuthService) Authenticate(token string) (string, string, error) {
	req, err := http.NewRequest("POST", s.railsBaseURL+"/verify_token", nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf(errResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authResp struct {
		UserID json.Number `json:"user_id"`
		Role   string      `json:"role"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", "", err
	}

	userID := string(authResp.UserID)
	return userID, strings.ToLower(authResp.Role), nil
}
