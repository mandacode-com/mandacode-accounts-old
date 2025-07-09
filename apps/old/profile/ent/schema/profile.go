package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Profile holds the schema definition for the Profile entity.
type Profile struct {
	ent.Schema
}

// Fields of the Profile.
func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Comment("The unique identifier for the profile"),
		field.String("email").
			Optional().
			Comment("The email address associated with the profile"),
		field.String("display_name").
			Optional().
			Comment("The display name of the profile owner"),
		field.String("bio").
			Optional().
			Comment("A short biography of the profile owner"),
		field.String("avatar_url").
			Optional().
			Comment("URL to the profile's avatar image"),
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

// Edges of the Profile.
func (Profile) Edges() []ent.Edge {
	return nil
}
