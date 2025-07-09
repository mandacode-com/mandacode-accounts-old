package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// LocalAuth holds the schema definition for the LocalAuth entity.
type LocalAuth struct {
	ent.Schema
}

// Fields of the LocalAuth.
func (LocalAuth) Fields() []ent.Field {
	return []ent.Field{
		// Internal PK
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			Default(uuid.New).
			Comment("The unique identifier for the local authentication record"),

		// AuthAccountID
		field.UUID("auth_account_id", uuid.UUID{}).
			Unique().
			Comment("The unique identifier for the authentication account associated with this local authentication"),

		// Email
		field.String("email").
			NotEmpty().
			Unique().
			Comment("The email address associated with the local authentication"),

		// Password
		field.String("password").
			NotEmpty().
			Comment("The hashed password for the local authentication"),

		// IsActive
		field.Bool("is_active").
			Default(true).
			Comment("Indicates if the local authentication is active and can be used to log in"),

		// IsVerified
		field.Bool("is_verified").
			Default(false).
			Comment("Indicates if the local authentication has verified the email address"),

		// CreatedAt
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the local authentication was created"),

		// UpdatedAt
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the local authentication was last updated"),

		// LastLoginAt
		field.Time("last_login_at").
			Optional().
			Comment("The time when the user last logged in with this local authentication"),

		// LastFailedLoginAt
		field.Time("last_failed_login_at").
			Optional().
			Comment("The time when the user last failed to log in with this local authentication"),

		// FailedLoginAttempts
		field.Int("failed_login_attempts").
			Default(0).
			Comment("The number of consecutive failed login attempts for this local authentication"),
	}
}

// Indexes of the LocalAuth.
func (LocalAuth) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("auth_account_id").Unique(),
	}
}

// Edges of the LocalAuth.
func (LocalAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("auth_account", AuthAccount.Type).
			Ref("local_auths").
			Unique().
			Field("auth_account_id").
			Required().
			Comment("The authentication account associated with this local authentication"),
	}
}
