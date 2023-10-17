//go:build ignore

// See Upstream docs for more details: https://entgo.io/docs/code-gen/#use-entc-as-a-package

package main

import (
	"log"
	"os"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/hedwigz/entviz"
)

func main() {
	// Ensure the schema directory exists before running entc.
	_ = os.Mkdir("schema", 0755)

	gqlExt, err := entgql.NewExtension(
		// Tell Ent to generate a GraphQL schema for
		// the Ent schema in a file named ent.graphql.
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("schema/ent.graphql"),
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithWhereInputs(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	opts := []entc.Option{
		entc.Extensions(
			gqlExt,
			entviz.Extension{}, // graph visualization
		),
	}

	if err := entc.Generate("./internal/ent/schema", &gen.Config{
		Target:   "./internal/ent/generated",
		Features: []gen.Feature{gen.FeatureVersionedMigration},
		Package:  "github.com/datumforge/go-template/internal/ent/generated",
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
