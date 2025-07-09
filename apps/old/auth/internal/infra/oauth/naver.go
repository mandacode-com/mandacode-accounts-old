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

// NaverOAuthProvider implements the OAuthService interface for Naver.
type NaverOAuthProvider struct {
	userInfoEndpoint string
	validate         *validator.Validate
}

// RawNaverUserInfo represents the raw user info structure returned by Naver OAuth.
type RawNaverUserInfo struct {
	ID         string `json:"response.id"`
	Email      string `json:"response.email"`
	Name       string `json:"response.nickname"`
	ResultCode string `json:"resultcode"`
}

func NewNaverOAuthProvider(
	validate *validator.Validate,
) oauthdomain.OAuthProvider {
	return &NaverOAuthProvider{
		userInfoEndpoint: oauthmeta.NaverUserInfoEndpoint,
		validate:         validate,
	}
}

func (s *NaverOAuthProvider) GetUserInfo(accessToken string) (*oauthdto.OAuthUserInfo, error) {
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

	// Decode the JSON response into the RawNaverUserInfo struct
	var rawUserInfo RawNaverUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	if rawUserInfo.ResultCode != "00" {
		return nil, fmt.Errorf("failed to fetch user info: result code %s", rawUserInfo.ResultCode)
	}

	oauthUserInfo := oauthdto.NewOAuthUserInfo(
		rawUserInfo.ID,
		rawUserInfo.Email,
		rawUserInfo.Name,
		true, // Naver does not provide email verification status, assuming true
	)
	if err := s.validate.Struct(oauthUserInfo); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return oauthUserInfo, nil
}
