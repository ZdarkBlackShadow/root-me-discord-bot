package services

import (
	"encoding/json"
	"os"
	"rootme-bot/models"
	"rootme-bot/repository"
	"strconv"

	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo: repo,
		logger: logger,
	}
}

func (s *UserService) SyncUsersFromJson(filePath string) error {
	s.logger.Info("User synchronization", zap.String("file", filePath))

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		s.logger.Error("Error while trying to reading JSON file", zap.Error(err))
		return err
	}

	var usersJson []models.UserJson
	if err := json.Unmarshal(fileData, &usersJson); err != nil {
		s.logger.Error("JSON decoding error", zap.Error(err))
		return err
	}

	for _, uj := range usersJson {
		if uj.RootMeID == 0 || uj.DiscordName == "" {
			s.logger.Warn("Invalid JSON user ignored", zap.String("name", uj.DiscordName))
			continue
		}

		rootMeIDStr := strconv.Itoa(uj.RootMeID)

		userModel := &models.User{
			DiscordName: uj.DiscordName,
			RootMeID:    rootMeIDStr,
		}

		if err := s.repo.CreateOrUpdate(userModel); err != nil {
			s.logger.Error("Persistence error", 
				zap.String("user", uj.DiscordName), 
				zap.Error(err))
			continue
		}
	}

	s.logger.Info("Synchronization complete")
	return nil
}