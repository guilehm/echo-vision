//go:build exclude

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	err := entc.Generate("./internal/infra/postgres/schema", &gen.Config{
		Target:  "./internal/infra/postgres/generated/ent",
		Schema:  "github.com/guilehm/echo-vision/internal/infra/postgres/schema",
		Package: "github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent",
		Features: []gen.Feature{
			gen.FeatureLock,
			gen.FeatureVersionedMigration,
		},
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
