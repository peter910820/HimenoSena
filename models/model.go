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
	BotChannelID2 string
	VoiceManageID string
	DevCategoryID string
}
