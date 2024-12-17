package models

import (
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Bot           *discordgo.Session
	Token         string
	AppId         string
	VoiceManageId string
}
