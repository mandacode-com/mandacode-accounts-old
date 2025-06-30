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

// KakaoOAuthProvider implements the OAuthService interface for Kakao.
type KakaoOAuthProvider struct {
	userInfoEndpoint string
	validate         *validator.Validate
}

// RawKakaoUserInfo represents the raw user info structure returned by Kakao OAuth.
type RawKakaoUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"kakao_account.email"`
	Name          string `json:"properties.nickname"`
	EmailValid    bool   `json:"kakao_account.is_email_valid"`
	EmailVerified bool   `json:"kakao_account.is_email_verified"`
}

func NewKakaoOAuthProvider(
	validate *validator.Validate,
) oauthdomain.OAuthProvider {
	return &KakaoOAuthProvider{
		userInfoEndpoint: oauthmeta.KakaoUserInfoEndpoint,
		validate:         validate,
	}
}

func (s *KakaoOAuthProvider) GetUserInfo(accessToken string) (*oauthdto.OAuthUserInfo, error) {
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

	oauthUserInfo := oauthdto.NewOAuthUserInfo(
		rawUserInfo.ID,
		rawUserInfo.Email,
		rawUserInfo.Name,
		rawUserInfo.EmailValid && rawUserInfo.EmailVerified,
	)
	if err := s.validate.Struct(oauthUserInfo); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return oauthUserInfo, nil
}
