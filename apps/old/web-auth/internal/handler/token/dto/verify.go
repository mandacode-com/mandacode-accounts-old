package tokenhandlerdto

import "github.com/google/uuid"

type VerifyAccessTokenResponse struct {
	Valid  bool       `json:"valid"`
	UserID *uuid.UUID `json:"user_id,omitempty"`
}
