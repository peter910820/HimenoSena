package commands

import (
	"HimenoSena/utils"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate, delay time.Duration) {
	utils.SendInteractionMsg(s, i, fmt.Sprintf("延遲時間為: %v", delay))
}

func Send(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	utils.SendInteractionMsg(s, i, message)
}

func GetRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleMap := make(map[string]string)
	roleMap["galgame"] = os.Getenv("GALGAME_ROLES_ID")

	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "roles" {
			selected := opt.StringValue()
			val, ok := roleMap[selected]
			if ok {
				err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, val)
				if err != nil {
					utils.SendInteractionMsg(s, i, "")
				} else {
					utils.SendInteractionMsg(s, i, "身分組獲取成功！")
				}
			} else {
				utils.SendInteractionMsg(s, i, "")
			}
		}
	}
}
