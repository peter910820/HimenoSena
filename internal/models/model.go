package models

import (
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Bot           *discordgo.Session
	Token         string
	AppID         string
	MainGuildID   string
	BotChannelID  string
	VoiceManageID string
}
