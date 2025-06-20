package domain

// TokenGenerator defines an interface for generating JWT tokens
type TokenGenerator interface {
	GenerateToken(
		claims map[string]string,
	) (string, int64, error)
	VerifyToken(
		token string,
	) (map[string]string, error)
}
