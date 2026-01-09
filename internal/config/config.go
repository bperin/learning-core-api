package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DBURL           string
	GoogleProjectID string
	PubSubTopicID   string
	GoogleAPIKey    string
	JWTSecret       string
	FileStoreName   string
	GCSBucketName   string
	SignedURLTTL    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:            os.Getenv("PORT"),
		DBURL:           os.Getenv("DB_URL"),
		GoogleProjectID: os.Getenv("GOOGLE_PROJECT_ID"),
		PubSubTopicID:   os.Getenv("PUBSUB_TOPIC_ID"),
		GoogleAPIKey:    os.Getenv("GOOGLE_API_KEY"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		FileStoreName:   os.Getenv("FILE_STORE_NAME"),
		GCSBucketName:   os.Getenv("GCS_BUCKET_NAME"),
		SignedURLTTL:    os.Getenv("GCS_SIGNED_URL_TTL"),
	}
}
