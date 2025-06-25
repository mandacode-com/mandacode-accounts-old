package oauthsvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mandacode.com/accounts/auth/internal/domain/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	"mandacode.com/accounts/auth/internal/infra/oauth/endpoint"
)

// GoogleOAuthService implements the GoogleOAuthService interface.
type GoogleOAuthService struct {
	userInfoEndpoint string
}

// RawGoogleUserInfo represents the raw user info structure returned by Google OAuth.
type RawGoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
}

func NewGoogleOAuthService() oauthdomain.OAuthService {
	return &GoogleOAuthService{
		userInfoEndpoint: endpoint.GoogleUserInfoEndpoint,
	}
}

func (s *GoogleOAuthService) GetUserInfo(accessToken string) (*dto.OAuthUserInfo, error) {
	req, err := http.NewRequest("GET", s.userInfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: status code %d", resp.StatusCode)
	}

	// Decode the JSON response into the RawGoogleUserInfo struct
	var rawUserInfo RawGoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	if rawUserInfo.Sub == "" {
		return nil, fmt.Errorf("user info does not contain a valid provider ID")
	}

	if rawUserInfo.Email == "" {
		return nil, fmt.Errorf("user info does not contain a valid email")
	}

	if rawUserInfo.Name == "" {
		return nil, fmt.Errorf("user info does not contain a valid name")
	}

	if !rawUserInfo.EmailVerified {
		return nil, fmt.Errorf("user info indicates email is not verified")
	}

	return &dto.OAuthUserInfo{
		ProviderID:    rawUserInfo.Sub,
		Email:         rawUserInfo.Email,
		Name:          rawUserInfo.Name,
		EmailVerified: rawUserInfo.EmailVerified,
	}, nil
}
