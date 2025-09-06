package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena/bot"
	"HimenoSena/model"
	"HimenoSena/utils"
)

func GetLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, serverUserExp *model.ServerMemberExp) {
	serverUserExp.Mu.Lock()
	defer serverUserExp.Mu.Unlock()
	memberData, err := bot.QueryUser(i.Member.User.ID, db)
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
