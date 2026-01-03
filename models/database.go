package models

import (
	"time"
)

type User struct {
	ID            uint       `gorm:"primaryKey" json:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DiscordName   string    `gorm:"size:100"`
	RootMeID      string    `gorm:"uniqueIndex;not null"`
	LastScore     int       `gorm:"default:0"`
	LastChallenge string    `gorm:"size:255"`
	Rank          string    `gorm:"size:50"`
}
