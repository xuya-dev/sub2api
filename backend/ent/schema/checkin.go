package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Checkin struct {
	ent.Schema
}

func (Checkin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "checkins"},
	}
}

func (Checkin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Time("checkin_date").
			SchemaType(map[string]string{
				dialect.Postgres: "date",
			}),
		field.Float("reward_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Int("streak_days").
			Default(1),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}),
	}
}

func (Checkin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("checkins").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Checkin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "checkin_date").Unique(),
		index.Fields("user_id"),
	}
}
