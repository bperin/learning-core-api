package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"learning-core-api/internal/config"
	"learning-core-api/internal/gcp"
	"learning-core-api/internal/infra"
	"learning-core-api/internal/persistance/seeds"
	"learning-core-api/internal/persistance/store"
)

// @title Learning API
// @version 1.0
// @description API with password OAuth2 authentication
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// =======================
// OAuth2 (GLOBAL)
// =======================

// @securityDefinitions.oauth2.password OAuth2
// @tokenUrl /oauth/token
func main() {
	// 1. Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Load Config
	cfg := config.Load()
	logConfig(cfg)

	// 3. Connect to Database
	db, err := infra.ConnectDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Fatal: could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// 3.5. Seed database if needed
	queries := store.New(db)
	if err := seeds.RunWithQueries(ctx, queries); err != nil {
		log.Printf("Warning: failed to seed database: %v", err)
	}

	// 4. Init PubSub
	pubsubService, err := infra.NewPubSubService(ctx, cfg.GoogleProjectID, cfg.PubSubTopicID)
	if err != nil {
		log.Printf("Warning: could not initialize pubsub: %v", err)
	} else {
		defer pubsubService.Close()
		log.Println("PubSub service initialized")
	}

	// 5. Init GCS
	var gcsService *gcp.GCSService
	if cfg.GCSBucketName != "" && cfg.FileStoreName != "" {
		var err error
		gcsService, err = gcp.NewGCSServiceFromConfig(ctx, cfg)
		if err != nil {
			log.Printf("Warning: could not initialize gcs service: %v", err)
		} else {
			defer gcsService.Close()
			_, err := gcp.NewFileServiceFromConfig(ctx, cfg, gcsService)
			if err != nil {
				log.Printf("Warning: could not initialize file service for store %q: %v", cfg.FileStoreName, err)
			} else {
				log.Println("File service initialized")
			}
		}
	} else {
		log.Println("File service not initialized (missing bucket or store name)")
	}

	// 6. Start HTTP Server
	router := infra.NewRouter(infra.RouterDeps{
		JWTSecret:    cfg.JWTSecret,
		Queries:      queries,
		GoogleAPIKey: cfg.GoogleAPIKey,
		GCSService:   gcsService,
	})
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func logConfig(cfg *config.Config) {
	if cfg == nil {
		log.Println("Config: <nil>")
		return
	}
	log.Printf(
		"Config loaded: port=%s db_url=%s project_id=%s pubsub_topic=%s file_store_name=%s gcs_bucket=%s signed_url_ttl=%s google_api_key=%s jwt_secret=%s",
		cfg.Port,
		redact(cfg.DBURL),
		cfg.GoogleProjectID,
		cfg.PubSubTopicID,
		cfg.FileStoreName,
		cfg.GCSBucketName,
		cfg.SignedURLTTL,
		redact(cfg.GoogleAPIKey),
		redact(cfg.JWTSecret),
	)

	if creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); creds != "" {
		log.Printf("GOOGLE_APPLICATION_CREDENTIALS=%s", creds)
	}
}

func redact(value string) string {
	if value == "" {
		return "<empty>"
	}
	if len(value) <= 4 {
		return "***"
	}
	return "***" + value[len(value)-4:]
}
