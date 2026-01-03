package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rootme-bot/models"
	"rootme-bot/repository"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type RootMeService struct {
	repo           *repository.UserRepository
	logger         *zap.Logger
	httpClient     *http.Client
	apiKey         string
	discordService *DiscordService
}

func NewRootMeService(repo *repository.UserRepository, logger *zap.Logger, api_key string, discord_service *DiscordService) *RootMeService {
	return &RootMeService{
		repo:   repo,
		logger: logger,
		apiKey: api_key,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		discordService: discord_service,
	}
}

// CheckUpdates parcourt tous les utilisateurs et vérifie les nouveaux flags
func (s *RootMeService) CheckUpdates() {
	s.logger.Info("Starting verification cycle")

	users, err := s.repo.GetAllUsers()
	if err != nil {
		s.logger.Error("Unable to retrieve users", zap.Error(err))
		return
	}

	for _, user := range users {
		s.logger.Debug("User verification", zap.String("pseudo", user.DiscordName))

		profile, err := s.FetchRootMeProfile(user.RootMeID)
		if err != nil {
			s.logger.Warn("Profile recovery failed", zap.String("id", user.RootMeID), zap.Error(err))
			continue
		}

		s.logger.Debug("profile : ", zap.Any("struct", profile))
		currentScore, err := strconv.Atoi(profile.Score)
		if err != nil {
			s.logger.Error("error when trying to convert str to int", zap.Error(err))
			return
		}

		s.logger.Debug("last score and current score", zap.Int("last score", user.LastScore), zap.Int("current score", currentScore))
		if user.LastScore != 0 && currentScore > user.LastScore {
			s.handleProgression(user, profile, currentScore)
		} else if user.LastScore == 0 {
			s.logger.Info("Initialization of the score in the database", zap.String("user", user.DiscordName), zap.Int("score", currentScore))
			s.repo.UpdateScore(user.RootMeID, currentScore, "")
		}

		time.Sleep(3 * time.Second) //for the rate limiting
	}
}

func (s *RootMeService) FetchRootMeProfile(rootMeID string) (*models.RootMeProfile, error) {
	url := fmt.Sprintf("https://api.www.root-me.org/auteurs/%s", rootMeID)

	// 1. Création de la requête
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	cookie := &http.Cookie{
		Name:  "api_key",
		Value: s.apiKey,
	}
	req.AddCookie(cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("root me API error: status %d", resp.StatusCode)
	}

	var profile models.RootMeProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *RootMeService) handleProgression(user models.User, profile *models.RootMeProfile, newScore int) {
	lastChallengeTitle := "Unknown challenge"
	if len(profile.Validations) > 0 {
		lastChallengeTitle = profile.Validations[0].Titre //the first challenge on the list is the last one completed
	}

	pointsGained := newScore - user.LastScore

	s.logger.Info("New challenge completed !",
		zap.String("user", user.DiscordName),
		zap.String("challenge", lastChallengeTitle),
		zap.Int("points", pointsGained),
	)

	err := s.repo.UpdateScore(user.RootMeID, newScore, lastChallengeTitle)
	if err != nil {
		s.logger.Error("error when trying to update score in database", zap.Error(err))
	}

	if s.discordService != nil {
		notif := models.DiscordNotification{
			Username:      user.DiscordName,
			ChallengeName: lastChallengeTitle,
			PointsGained:  pointsGained,
			TotalScore:    newScore,
			Position:      profile.Position,
			RootMeID:      user.RootMeID,
		}

		err := s.discordService.SendProgressionEmbed(notif)
		if err != nil {
			s.logger.Error("Failed to send Discord notification", zap.Error(err))
		}
	}
}
