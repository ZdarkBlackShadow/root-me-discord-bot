package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	DB           *gorm.DB
	Logger       *zap.Logger
	DiscordToken string
	ChannelID    string
	RootMeApiKey string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	logger := initLogger()

	db , err := initDB()
	if err != nil {
		logger.Error("Problème when trying to connect to database : ", zap.Error(err))
	}

	return &Config{
		DB:           db,
		Logger:       logger,
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		ChannelID:    os.Getenv("DISCORD_CHANNEL_ID"),
		RootMeApiKey: os.Getenv("API_KEY"),
	}
}
