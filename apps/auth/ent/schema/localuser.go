package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LocalUser holds the schema definition for the LocalUser entity.
type LocalUser struct {
	ent.Schema
}

// Annotations of the LocalUser.
func (LocalUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Schema("public"),
	}
}

// Fields of the LocalUser.
func (LocalUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique(),
		field.String("email").Unique().NotEmpty().Unique(),
		field.String("password").NotEmpty(),
		field.Bool("is_active").Default(true),
		field.Bool("is_verified").Default(false),
	}
}

// Edges of the LocalUser.
func (LocalUser) Edges() []ent.Edge {
	return nil
}
