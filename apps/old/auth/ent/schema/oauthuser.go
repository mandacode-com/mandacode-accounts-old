package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// OAuthUser holds the schema definition for the OAuthUser entity.
type OAuthUser struct {
	ent.Schema
}

// Annotations of the OAuthUser.
func (OAuthUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Schema("public"),
	}
}

// Fields of the OAuthUser.
func (OAuthUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Comment("The unique identifier for the OAuth user"),
		field.String("email").
			NotEmpty().
			Comment("The email address associated with the OAuth user"),
		field.Enum("provider").
			Values("google", "github", "facebook", "kakao", "naver", "apple").
			Comment("The OAuth provider used for authentication"),
		field.String("provider_id").
			NotEmpty().
			Comment("The unique identifier provided by the OAuth provider for the user"),
		field.Bool("is_active").
			Default(true).
			Comment("Indicates if the OAuth user is active and can log in"),
		field.Bool("is_verified").
			Default(true).
			Comment("Indicates if the OAuth user has verified their email address"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the OAuth user was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the OAuth user was last updated"),
	}
}

// Indexes of the OAuthUser.
func (OAuthUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "provider").Unique(),
		index.Fields("provider", "provider_id").Unique(),
	}
}

// Edges of the OAuthUser.
func (OAuthUser) Edges() []ent.Edge {
	return nil
}
