package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// OAuthAuth holds the schema definition for the OAuthAuth entity.
type OAuthAuth struct {
	ent.Schema
}

// Fields of the OAuthAuth.
func (OAuthAuth) Fields() []ent.Field {
	return []ent.Field{
		// Internal PK
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			Default(uuid.New).
			Comment("The unique identifier for the OAuth authentication record"),

		// AuthAccountID
		field.UUID("auth_account_id", uuid.UUID{}).
			Comment("The unique identifier for the authentication account associated with this OAuth authentication"),

		// Provider
		field.Enum("provider").
			Values("google", "github", "facebook", "kakao", "naver", "apple").
			Comment("The OAuth provider used for authentication"),

		// ProviderID
		field.String("provider_id").
			NotEmpty().
			Comment("The unique identifier provided by the OAuth provider for the user"),

		// IsActive
		field.Bool("is_active").
			Default(true).
			Comment("Indicates if the OAuth authentication is active and can be used to log in"),

		// IsVerified
		field.Bool("is_verified").
			Default(true).
			Comment("Indicates if the OAuth authentication has verified the email address"),

		// CreatedAt
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the OAuth authentication was created"),

		// UpdatedAt
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the OAuth authentication was last updated"),

		// Email
		field.String("email").
			Optional().
			Comment("The email address associated with the OAuth authentication"),

		// LastLoginAt
		field.Time("last_login_at").
			Optional().
			Comment("The time when the user last logged in with this OAuth authentication"),

		// LastFailedLoginAt
		field.Time("last_failed_login_at").
			Optional().
			Comment("The time when the user last failed to log in with this OAuth authentication"),

		// FailedLoginAttempts
		field.Int("failed_login_attempts").
			Default(0).
			Comment("The number of consecutive failed login attempts for this OAuth authentication"),
	}
}

// Indexes of the OAuthAuth.
func (OAuthAuth) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("auth_account_id", "provider").Unique(),
		index.Fields("provider", "provider_id").Unique(),
	}
}

// Edges of the OAuthAuth.
func (OAuthAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("auth_account", AuthAccount.Type).
			Ref("oauth_auths").
			Unique().
			Required().
			Field("auth_account_id").
			Comment("The authentication account associated with this OAuth authentication"),
	}
}
