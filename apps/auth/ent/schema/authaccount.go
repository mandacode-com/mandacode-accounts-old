package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AuthAccount holds the schema definition for the AuthAccount entity.
type AuthAccount struct {
	ent.Schema
}

// Fields of the AuthAccount.
func (AuthAccount) Fields() []ent.Field {
	return []ent.Field{
		// Internal PK
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			Default(uuid.New).
			Comment("The unique identifier for the authentication account"),

		// User ID
		field.UUID("user_id", uuid.UUID{}).
			Comment("The unique identifier for the user associated with this authentication account"),

		// CreatedAt
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the authentication account was created"),

		// UpdatedAt
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the authentication account was last updated"),

		// LastLoginAt
		field.Time("last_login_at").
			Optional().
			Comment("The time when the user last logged in with this authentication account"),

		// LastFailedLoginAt
		field.Time("last_failed_login_at").
			Optional().
			Comment("The time when the user last failed to log in with this authentication account"),

		// FailedLoginAttempts
		field.Int("failed_login_attempts").
			Default(0).
			Comment("The number of consecutive failed login attempts for this authentication account"),
	}
}

// Edges of the AuthAccount.
func (AuthAccount) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge to LocalAuth auth_account_id
		edge.To("local_auths", LocalAuth.Type).
			Comment("The local authentication methods associated with this authentication account"),

		// Edge to OAuthAuth auth_account_id
		edge.To("oauth_auths", OAuthAuth.Type).
			Comment("The OAuth authentication methods associated with this authentication account"),
	}
}
