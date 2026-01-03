package models

type DiscordNotification struct {
	Username      string
	ChallengeName string
	PointsGained  int
	TotalScore    int
	Position      int
	LogoURL       string
	RootMeID      string
}
