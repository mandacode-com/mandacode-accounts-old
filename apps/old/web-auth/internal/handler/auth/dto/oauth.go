package authhandlerdto

type OAuthLoginURLResponse struct {
	URL string `json:"url"`
}

type OAuthCallbackResponse struct {
	AccessToken string `json:"access_token"`
	// RefreshToken will be saved in the session
}
