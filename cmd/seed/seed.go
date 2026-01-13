package main

import (
	"context"

	"learning-core-api/internal/persistance/seeds"
	"learning-core-api/internal/persistance/store"
)

func runSeeds(queries *store.Queries) error {
	return seeds.RunWithQueries(context.Background(), queries)
}
