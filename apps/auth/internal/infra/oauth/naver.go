package oauthsvc

import (
	"encoding/json"
	"fmt"
	"mandacode.com/accounts/auth/internal/domain/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	"mandacode.com/accounts/auth/internal/infra/oauth/endpoint"
	"net/http"
)

// NaverOAuthService implements the OAuthService interface for Naver.
type NaverOAuthService struct {
	userInfoEndpoint string
}

// RawNaverUserInfo represents the raw user info structure returned by Naver OAuth.
type RawNaverUserInfo struct {
	ID         string `json:"response.id"`
	Email      string `json:"response.email"`
	Name       string `json:"response.nickname"`
	ResultCode string `json:"resultcode"`
}

func NewNaverOAuthService() oauthdomain.OAuthService {
	return &NaverOAuthService{
		userInfoEndpoint: endpoint.NaverUserInfoEndpoint,
	}
}

func (s *NaverOAuthService) GetUserInfo(accessToken string) (*dto.OAuthUserInfo, error) {
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

	if rawUserInfo.ID == "" {
		return nil, fmt.Errorf("user info does not contain a valid provider ID")
	}

	if rawUserInfo.Email == "" {
		return nil, fmt.Errorf("user info does not contain a valid email")
	}

	if rawUserInfo.Name == "" {
		return nil, fmt.Errorf("user info does not contain a valid name")
	}

	oauthUserInfo := dto.NewOAuthUserInfo(
		rawUserInfo.ID,
		rawUserInfo.Email,
		rawUserInfo.Name,
		true, // Naver does not provide email verification status, assuming true
	)
	return oauthUserInfo, nil
}
