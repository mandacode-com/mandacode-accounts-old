package authhandlerdto

type LocalLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LocalLoginResponse struct {
	Code string `json:"code"`
}

type LocalVerifyCodeResponse struct {
	AccessToken string `json:"access_token"`
	// RefreshToken will be saved in the session
}
