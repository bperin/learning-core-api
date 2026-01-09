package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DBURL             string
	GoogleProjectID   string
	PubSubTopicID     string
	GoogleAPIKey      string
	JWTSecret         string
	LearningStoreName string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:              os.Getenv("PORT"),
		DBURL:             os.Getenv("DB_URL"),
		GoogleProjectID:   os.Getenv("GOOGLE_PROJECT_ID"),
		PubSubTopicID:     os.Getenv("PUBSUB_TOPIC_ID"),
		GoogleAPIKey:      os.Getenv("GOOGLE_API_KEY"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		LearningStoreName: os.Getenv("LEARNING_STORE_NAME"),
	}
}
