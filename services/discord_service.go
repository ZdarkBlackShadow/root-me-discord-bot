package services

import (
	"fmt"
	"rootme-bot/models"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type DiscordService struct {
	session   *discordgo.Session
	channelID string
	logger    *zap.Logger
}

func NewDiscordService(token string, channelID string, logger *zap.Logger) (*DiscordService, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordService{
		session:   dg,
		channelID: channelID,
		logger:    logger,
	}, nil
}

func (s *DiscordService) Open() error {
	return s.session.Open()
}

func (s *DiscordService) Close() {
	s.session.Close()
}

func (s *DiscordService) SendProgressionEmbed(notif models.DiscordNotification) error {
	embed := &discordgo.MessageEmbed{
		Title:       "🚩 New Challenge Validated !",
		Description: fmt.Sprintf("**%s** has passed another challenge !", notif.Username),
		Color:       0x47b447,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Challenge",
				Value:  fmt.Sprintf("`%s`", notif.ChallengeName),
				Inline: true,
			},
			{
				Name:   "Points earned",
				Value:  fmt.Sprintf("+%d pts", notif.PointsGained),
				Inline: true,
			},
			{
				Name:   "New Score",
				Value:  fmt.Sprintf("%d", notif.TotalScore),
				Inline: false,
			},
			{
				Name: "Ranking",
				Value: fmt.Sprintf("%d", notif.Position),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Root-Me Bot",
		},
		URL: fmt.Sprintf("https://www.root-me.org/%s", notif.Username),
	}

	_, err := s.session.ChannelMessageSendEmbed(s.channelID, embed)
	if err != nil {
		s.logger.Error("Error when send embed Discord", zap.Error(err))
		return err
	}

	s.logger.Info("Discord Notification sent", zap.String("user", notif.Username))
	return nil
}
