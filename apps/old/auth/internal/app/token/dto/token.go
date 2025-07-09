package tokendto

type NewToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyTokenResult struct {
	UserID string `json:"user_id"`
	Valid  bool   `json:"valid"`
}
