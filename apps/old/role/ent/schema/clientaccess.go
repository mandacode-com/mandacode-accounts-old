package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ClientAccess holds the schema definition for the ClientAccess entity.
type ClientAccess struct {
	ent.Schema
}

// Fields of the ServiceClient.
func (ClientAccess) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique(),
		field.UUID("service_id", uuid.UUID{}),
		field.String("name").
			NotEmpty(),
		field.String("client_id").
			Unique().
			NotEmpty(),
		field.String("client_secret").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.Bool("is_active").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the ServiceClient.
func (ClientAccess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).
			Ref("client_accesses").
			Field("service_id").
			Unique().
			Required(),
	}
}
