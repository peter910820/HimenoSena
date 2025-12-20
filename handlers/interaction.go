package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena"
	"HimenoSena/commands"
)

func OnInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, serverMemberExp *HimenoSena.ServerMemberExp, c *HimenoSena.Config) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		delay := s.HeartbeatLatency()
		go commands.Ping(s, i, delay)
	case "send":
		amount := i.ApplicationCommandData().Options[0].StringValue()
		go commands.Send(s, i, amount)
	case "取得身分組":
		go commands.GetRoles(s, i)
	case "取得聊天等級":
		go commands.GetLevel(s, i, db, serverMemberExp)
	case "取得群組等級排行":
		go commands.GetAllLevel(s, i, db, c)
	case "查詢log":
		go commands.GetLog(s, i)
	default:
		logrus.Warn("command not founds")
	}
}
