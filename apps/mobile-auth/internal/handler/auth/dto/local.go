package authhandlerdto

type LocalLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LocalLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
