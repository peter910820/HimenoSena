package bot

import (
	"os"

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
		{
			Name:        "取得聊天等級",
			Description: "取得自己的聊天等級以及經驗值",
		},
		{
			Name:        "取得群組等級排行",
			Description: "取得當前群組等級排行",
		},
	}

	// 私有指令，要使用群組內部整合管理複寫權限，預設是全部可見
	privateCommand := []*discordgo.ApplicationCommand{
		{
			Name:        "查詢log",
			Description: "查詢log用",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sites",
					Description: "選擇要查看的站台",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "kurohelper",
							Value: "kurohelper",
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

	guildID := os.Getenv("MAIN_GUILD_ID")
	guildID2 := os.Getenv("LOG_GUILD_ID")
	for _, cmd := range privateCommand {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			logrus.Error(err)
		}
		_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID2, cmd)
		if err != nil {
			logrus.Error(err)
		}
	}
}
