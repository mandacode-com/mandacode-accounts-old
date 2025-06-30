package oauthprovider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	oauthdto "mandacode.com/accounts/auth/internal/infra/oauth/dto"
	oauthmeta "mandacode.com/accounts/auth/internal/infra/oauth/meta"
)

// GoogleOAuthProvider implements the GoogleOAuthProvider interface.
type GoogleOAuthProvider struct {
	userInfoEndpoint string
	validate         *validator.Validate
}

// RawGoogleUserInfo represents the raw user info structure returned by Google OAuth.
type RawGoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
}

func NewGoogleOAuthProvider(
	validate *validator.Validate,
) oauthdomain.OAuthProvider {
	return &GoogleOAuthProvider{
		userInfoEndpoint: oauthmeta.GoogleUserInfoEndpoint,
		validate:         validate,
	}
}

func (s *GoogleOAuthProvider) GetUserInfo(accessToken string) (*oauthdto.OAuthUserInfo, error) {
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

	oauthUserInfo := oauthdto.NewOAuthUserInfo(
		rawUserInfo.Sub,
		rawUserInfo.Email,
		rawUserInfo.Name,
		rawUserInfo.EmailVerified,
	)
	if err := s.validate.Struct(oauthUserInfo); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return oauthUserInfo, nil
}
