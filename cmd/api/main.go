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
	"learning-core-api/internal/infra"
	"learning-core-api/internal/persistance/store"
)

// @securityDefinitions.oauth2 OAuth2Auth
// @type oauth2
// @flow password
// @authorizationUrl https://localhost:8080/oauth/authorize
// @tokenUrl https://localhost:8080/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
func main() {
	// 1. Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Load Config
	cfg := config.Load()

	// 3. Connect to Database
	db, err := infra.ConnectDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Fatal: could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// 4. Init PubSub
	pubsubService, err := infra.NewPubSubService(ctx, cfg.GoogleProjectID, cfg.PubSubTopicID)
	if err != nil {
		log.Printf("Warning: could not initialize pubsub: %v", err)
	} else {
		defer pubsubService.Close()
		log.Println("PubSub service initialized")
	}

	// 5. Start HTTP Server
	queries := store.New(db)
	router := infra.NewRouter(infra.RouterDeps{
		JWTSecret:    cfg.JWTSecret,
		Queries:      queries,
		GoogleAPIKey: cfg.GoogleAPIKey,
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
