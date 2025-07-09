package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// GroupUser holds the schema definition for the GroupUser entity.
type GroupUser struct {
	ent.Schema
}

// Fields of the GroupUser.
func (GroupUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("group_id", uuid.UUID{}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the GroupUser.
func (GroupUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("group_users").
			Field("group_id").
			Unique().
			Required(),
	}
}

func (GroupUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("group_id", "user_id").
			Unique(),
	}
}
