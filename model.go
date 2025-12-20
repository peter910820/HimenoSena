package HimenoSena

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Bot            *discordgo.Session
	Token          string
	AppID          string
	MainGuildID    string
	BotChannelID   string
	BotChannelID2  string
	LevelUpChannel string
	VoiceManageID  string
	DevCategoryID  string
}

type ServerMemberExp struct {
	ServerID   string
	MemberData map[string]uint // key is userID, value is exp required to upgrade
	Mu         sync.Mutex
}
