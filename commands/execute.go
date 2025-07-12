package commands

import (
	"HimenoSena/models"
	"HimenoSena/utils"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func GetChatLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, serverUserExp *models.ServerMemberExp) {
	serverUserExp.Mu.Lock()
	defer serverUserExp.Mu.Unlock()
	memberData, err := utils.QueryUser(i.Member.User.ID, db)
	if err != nil {
		logrus.Error(err)
		return
	}
	val, ok := serverUserExp.MemberData[i.Member.User.ID]
	if !ok {
		logrus.Error("找不到該使用者的經驗值資料")
		return
	}
	utils.SendInteractionMsg(s, i, fmt.Sprintf("**%s** 目前等級為 **%d** 等，距離下一等還差 **%d** 經驗值", i.Member.User.GlobalName, memberData.Level, val))
}
