package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Todo holds the example schema definition for the Todo entity
type Todo struct {
	ent.Schema
}

// Fields of the Todo
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Immutable(),
		field.String("name").
			Comment("the name of the organization").
			NotEmpty(),
		field.String("description").
			Comment("An optional description of the organization").
			Optional(),
	}
}

func (Todo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

// Annotations of the Organization
func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
