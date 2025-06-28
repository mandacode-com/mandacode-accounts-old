package oauthsvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mandacode.com/accounts/auth/internal/domain/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	"mandacode.com/accounts/auth/internal/infra/service/oauth/endpoint"
)

// KakaoOAuthService implements the OAuthService interface for Kakao.
type KakaoOAuthService struct {
	userInfoEndpoint string
}

// RawKakaoUserInfo represents the raw user info structure returned by Kakao OAuth.
type RawKakaoUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"kakao_account.email"`
	Name          string `json:"properties.nickname"`
	EmailValid    bool   `json:"kakao_account.is_email_valid"`
	EmailVerified bool   `json:"kakao_account.is_email_verified"`
}

func NewKakaoOAuthService() oauthdomain.OAuthService {
	return &KakaoOAuthService{
		userInfoEndpoint: endpoint.KakaoUserInfoEndpoint,
	}
}

func (s *KakaoOAuthService) GetUserInfo(accessToken string) (*dto.OAuthUserInfo, error) {
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

	// Decode the JSON response into the RawKakaoUserInfo struct
	var rawUserInfo RawKakaoUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
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
		rawUserInfo.EmailValid && rawUserInfo.EmailVerified,
	)
	return oauthUserInfo, nil
}
