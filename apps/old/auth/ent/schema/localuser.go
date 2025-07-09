package schema

import (
	"time"

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
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique(),
		field.String("email").
			Unique().
			NotEmpty().
			Unique().
			Comment("The email address of the user, must be unique"),
		field.String("password").
			NotEmpty().
			Comment("The hashed password of the user"),
		field.Bool("is_active").Default(true).
			Comment("The user is active and can log in"),
		field.Bool("is_verified").
			Default(false).
			Comment("The user has verified their email address"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the profile was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the profile was last updated"),
	}
}

// Edges of the LocalUser.
func (LocalUser) Edges() []ent.Edge {
	return nil
}
