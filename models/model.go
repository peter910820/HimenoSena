package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Bot           *discordgo.Session
	Token         string
	AppID         string
	MainGuildID   string
	BotChannelID  string
	BotChannelID2 string
	VoiceManageID string
	DevCategoryID string
}

type User struct {
	UserID    string    `gorm:"primaryKey"`
	ServerID  string    `gorm:"primaryKey"`
	Level     uint      `gorm:"not null;default:1"`
	Exp       uint      `gorm:"not null;default:0"`
	JoinTime  time.Time `gorm:"not null"`
	UpdatedAt time.Time
}
