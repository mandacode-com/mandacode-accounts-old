package authresdto

// LoginResponse is the response structure for login operations.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
