package commands

import (
	"HimenoSena/utils"

	"github.com/bwmarrin/discordgo"
)

func Send(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	utils.SendInteractionMsg(s, i, message)
}
