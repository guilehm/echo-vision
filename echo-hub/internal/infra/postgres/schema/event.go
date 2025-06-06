package schema

import (
	"encoding/json"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.Enum("type").
			Values(hubevents.EventType("").StringValues()...),
		field.Enum("sub_type").
			Values(hubevents.EventSubType("").StringValues()...),
		field.Enum("status").
			Values(hubevents.EventStatus("").StringValues()...),
		field.JSON("result", json.RawMessage{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("events").
			Unique().
			Required().
			Field("user_id"),
		edge.From("file", File.Type).
			Ref("events").
			Unique(),
	}
}
