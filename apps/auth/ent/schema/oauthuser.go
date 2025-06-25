package schema

import (
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
		field.UUID("id", uuid.UUID{}).Immutable(),
		field.String("email").NotEmpty(),
		field.Enum("provider").Values("google", "github", "facebook", "kakao", "naver", "apple"),
		field.String("provider_id").NotEmpty(),
		field.Bool("is_active").Default(true),
		field.Bool("is_verified").Default(true),
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
