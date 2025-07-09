package oauthcode

import (
	"encoding/json"
	"errors"
	"net/http"

	oauthcodedomain "mandacode.com/accounts/web-auth/internal/domain/oauthcode"
	codeapidto "mandacode.com/accounts/web-auth/internal/infra/oauthcode/dto/api"
	oauthcodemeta "mandacode.com/accounts/web-auth/internal/infra/oauthcode/meta"
)

type GoogleOAuthCode struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	GrantType    string
}

func NewGoogleOAuthCode(clientID, clientSecret, redirectURL string) oauthcodedomain.OAuthCode {
	return &GoogleOAuthCode{
		Endpoint:     oauthcodemeta.GoogleTokenEndpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		GrantType:    oauthcodemeta.GoogleGrantType,
	}
}

func (g *GoogleOAuthCode) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", g.Endpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", g.ClientID)
	q.Add("client_secret", g.ClientSecret)
	q.Add("redirect_uri", g.RedirectURL)
	q.Add("grant_type", g.GrantType)
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

	var tokenResponse codeapidto.GoogleAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty in response")
	}

	return tokenResponse.AccessToken, nil
}

func (g *GoogleOAuthCode) GetLoginURL() (string, error) {
	if g.ClientID == "" || g.RedirectURL == "" {
		return "", errors.New("client ID and redirect URL must be set")
	}

	loginURL := oauthcodemeta.GoogleAuthEndpoint + "?client_id=" + g.ClientID +
		"&redirect_uri=" + g.RedirectURL +
		"&response_type=code&scope=email%20profile&access_type=offline"

	return loginURL, nil
}
