package tokengendomain

type TokenGenerator interface {
	// GenerateToken generates a signed JWT token with the provided claims
	//
	// Parameters:
	//   - claims: a map of additional claims to include in the token
	//
	// Returns:
	//   - token: the generated JWT token as a string
	//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time of the token
	//   - error: an error if token generation fails (e.g., if the private key is not initialized or signing fails)
	GenerateToken(
		claims map[string]string,
	) (string, int64, error)

	// VerifyToken verifies the provided JWT token and extracts the claims
	//
	// Parameters:
	//   - token: the JWT token to verify
	//
	// Returns:
	//   - claims: a map of claims extracted from the token if valid
	//   - error: an error if verification fails (e.g., if the public key is not initialized or parsing fails)
	VerifyToken(
		token string,
	) (map[string]string, error)
}
