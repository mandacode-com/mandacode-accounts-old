package oauthcode

import (
	"encoding/json"
	"errors"
	"net/http"

	oauthcodedomain "mandacode.com/accounts/web-auth/internal/domain/oauthcode"
	codeapidto "mandacode.com/accounts/web-auth/internal/infra/oauthcode/dto/api"
	oauthcodemeta "mandacode.com/accounts/web-auth/internal/infra/oauthcode/meta"
)

type NaverOAuthCode struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	GrantType    string
}

func NewNaverOAuthCode(clientID, clientSecret, redirectURI string) oauthcodedomain.OAuthCode {
	return &NaverOAuthCode{
		Endpoint:     oauthcodemeta.NaverTokenEndpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		GrantType:    oauthcodemeta.NaverGrantType,
	}
}

func (n *NaverOAuthCode) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", n.Endpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", n.ClientID)
	q.Add("client_secret", n.ClientSecret)
	q.Add("redirect_uri", n.RedirectURI)
	q.Add("grant_type", n.GrantType)
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

	var tokenResponse codeapidto.NaverAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty")
	}

	return tokenResponse.AccessToken, nil
}

func (n *NaverOAuthCode) GetLoginURL() (string, error) {
	if n.ClientID == "" || n.RedirectURI == "" {
		return "", errors.New("client ID and redirect URI must be set")
	}

	url := n.Endpoint + "?client_id=" + n.ClientID +
		"&response_type=code&redirect_uri=" + n.RedirectURI

	return url, nil
}
