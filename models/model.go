package models

import (
	"sync"
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

type ServerMemberExp struct {
	ServerID   string
	MemberData map[string]uint // key is userID, value is exp required to upgrade
	Mu         sync.Mutex
}

// database schema
type Member struct {
	UserID     string    `gorm:"primaryKey"`
	ServerID   string    `gorm:"primaryKey"`
	UserName   string    `gorm:"not null"`
	Level      uint      `gorm:"not null;default:1"`
	Exp        uint      `gorm:"not null;default:0"` // 該等級的經驗值，加上LevelUpExp才是該成員的所有經驗值
	LevelUpExp uint      `gorm:"not null;default:5"`
	JoinAt     time.Time `gorm:"not null"`
	UpdatedAt  time.Time
}

type Log struct {
	LogID     string `gorm:"primaryKey"`
	ServerID  string `gorm:"primaryKey"`
	Type      string `gorm:"not null"`
	Message   string `gorm:"not null"`
	UpdatedAt time.Time
}
