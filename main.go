package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"rootme-bot/config"
	"rootme-bot/repository"
	"rootme-bot/services"

	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	logger := cfg.Logger
	defer logger.Sync()

	logger.Info("Starting the Root-Me-Bot")

	userRepo := repository.NewUserRepository(cfg.DB)

	userService := services.NewUserService(userRepo, logger)
	discordService, err := services.NewDiscordService(cfg.DiscordToken, cfg.ChannelID, logger)
	if err != nil {
		logger.Fatal("Error initializing the Discord service", zap.Error(err))
	}
	rootMeService := services.NewRootMeService(userRepo, logger, cfg.RootMeApiKey, discordService)

	if err := userService.SyncUsersFromJson("users.json"); err != nil {
		logger.Warn("Initial synchronization failed (users.json), using existing DB data")
	}

	if err := discordService.Open(); err != nil {
		logger.Fatal("Unable to connect the bot to Discord", zap.Error(err))
	}
	defer discordService.Close()
	logger.Info("Discord bot connected and ready")

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		rootMeService.CheckUpdates()

		for range ticker.C {
			rootMeService.CheckUpdates()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	logger.Info("Stopping the bot")
}
