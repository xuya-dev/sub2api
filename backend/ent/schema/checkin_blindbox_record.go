package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type CheckinBlindboxRecord struct {
	ent.Schema
}

func (CheckinBlindboxRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "checkin_blindbox_records"},
	}
}

func (CheckinBlindboxRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Int64("prize_item_id"),
		field.String("prize_name").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "varchar(100)"}),
		field.String("rarity").
			Default("common").
			SchemaType(map[string]string{dialect.Postgres: "varchar(20)"}),
		field.String("reward_type").
			Default("balance").
			SchemaType(map[string]string{dialect.Postgres: "varchar(30)"}),
		field.Float("reward_value").
			Default(0).
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Int("streak_days").
			Default(0),
		field.String("reward_detail").
			Optional().
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (CheckinBlindboxRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("created_at"),
	}
}
