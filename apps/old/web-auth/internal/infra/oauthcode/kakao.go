package oauthcode

import (
	"encoding/json"
	"errors"
	"net/http"

	oauthcodedomain "mandacode.com/accounts/web-auth/internal/domain/oauthcode"
	codeapidto "mandacode.com/accounts/web-auth/internal/infra/oauthcode/dto/api"
	oauthcodemeta "mandacode.com/accounts/web-auth/internal/infra/oauthcode/meta"
)

type KakaoOAuthCode struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	GrantType    string
}

func NewKakaoOAuthCode(clientID, clientSecret, redirectURL string) oauthcodedomain.OAuthCode {
	return &KakaoOAuthCode{
		Endpoint:     oauthcodemeta.KakaoTokenEndpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		GrantType:    oauthcodemeta.KakaoGrantType,
	}
}

func (k *KakaoOAuthCode) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", k.Endpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", k.ClientID)
	q.Add("client_secret", k.ClientSecret)
	q.Add("redirect_uri", k.RedirectURL)
	q.Add("grant_type", k.GrantType)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to retrieve access token: " + resp.Status)
	}

	var tokenResponse codeapidto.KakaoAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty")
	}

	return tokenResponse.AccessToken, nil
}

func (k *KakaoOAuthCode) GetLoginURL() (string, error) {
	if k.ClientID == "" || k.RedirectURL == "" {
		return "", errors.New("client ID and redirect URL must be set")
	}

	loginURL := oauthcodemeta.KakaoAuthEndpoint + "?client_id=" + k.ClientID +
		"&redirect_uri=" + k.RedirectURL +
		"&response_type=code&scope=account_email%20profile_nickname"

	return loginURL, nil
}
