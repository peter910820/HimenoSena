package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func BasicCommand(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "return bot heartbeatlatency",
		},
		{
			Name:        "send",
			Description: "use bot to send text message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "your text message",
					Required:    true,
				},
			},
		},
		{
			Name:        "取得身分組",
			Description: "取得身分組",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "roles",
					Description: "choice a role",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Galgame玩家",
							Value: "galgame",
						},
					},
				},
			},
		},
	}
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
}
