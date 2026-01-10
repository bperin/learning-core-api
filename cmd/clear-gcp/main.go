package main

import (
	"context"
	"log"

	"learning-core-api/internal/config"
	"learning-core-api/internal/gcp"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	gcsService, err := gcp.NewGCSServiceFromConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init gcs service: %v", err)
	}
	defer gcsService.Close()

	fileService, err := gcp.NewFileServiceFromConfig(ctx, cfg, gcsService)
	if err != nil {
		log.Fatalf("failed to init file service: %v", err)
	}

	log.Printf("Clearing GCS bucket: %s", cfg.GCSBucketName)
	if err := gcsService.EmptyBucket(ctx); err != nil {
		log.Printf("error clearing bucket: %v", err)
	} else {
		log.Println("GCS bucket cleared successfully")
	}

	log.Println("Clearing Gemini File Search Stores...")
	if err := fileService.ClearAllStores(ctx); err != nil {
		log.Printf("error clearing stores: %v", err)
	} else {
		log.Println("Gemini File Search Stores cleared successfully")
	}
}
