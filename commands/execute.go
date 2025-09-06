package commands

import (
	"HimenoSena/model"
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

func GetChatLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, serverUserExp *model.ServerMemberExp) {
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

func GetGroupAllLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, c *model.Config) {
	var memberData []model.Member
	err := db.Select("user_name, level, join_at").Where("server_id = ?", c.MainGuildID).Order("exp DESC").Limit(10).Find(&memberData).Error
	if err != nil {
		logrus.Error(err)
	}

	resultStr := ""
	for i, v := range memberData {
		resultStr += fmt.Sprintf("%d. **%s** %d等 加入時間: %s\n", i, v.UserName, v.Level, v.JoinAt.Format("2006-01-02"))
	}

	utils.SendInteractionMsg(s, i, resultStr)
}
