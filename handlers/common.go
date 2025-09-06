package handlers

import (
	"HimenoSena/bot"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "恋×シンアイ彼女")
	bot.BasicCommand(s)
}
