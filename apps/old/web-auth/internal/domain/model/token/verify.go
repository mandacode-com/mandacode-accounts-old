package tokenmodel

import "github.com/google/uuid"

type VerifyAccessTokenResult struct {
	Valid  bool      `json:"valid"`
	UserID *uuid.UUID `json:"user_id,omitempty"`
}
