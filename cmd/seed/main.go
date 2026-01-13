package main

import (
	"log"

	"learning-core-api/internal/config"
	"learning-core-api/internal/infra"
	"learning-core-api/internal/persistance/store"
)

func main() {
	cfg := config.Load()
	if cfg.DBURL == "" {
		log.Fatal("DB_URL is required")
	}

	db, err := infra.ConnectDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	queries := store.New(db)
	if err := runSeeds(queries); err != nil {
		log.Fatalf("failed to seed: %v", err)
	}

	log.Println("seed complete")
}
