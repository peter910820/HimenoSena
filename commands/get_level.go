package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena/bot"
	"HimenoSena/models"
	"HimenoSena/utils"
)

func GetLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, serverUserExp *models.ServerMemberExp) {
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

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("🔔**%s等級資訊**", memberData.UserName),
		Color: 0xB5CAA0,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "目前等級/總經驗值",
				Value:  fmt.Sprintf("**%dLv**/%d", memberData.Level, memberData.Exp+memberData.LevelUpExp),
				Inline: true,
			},
			{
				Name:   "距離升級經驗",
				Value:  fmt.Sprintf("%d", val),
				Inline: true,
			},
			{
				Name:   "加入時間",
				Value:  memberData.JoinAt.Format("2006-01-02"),
				Inline: true,
			},
		},
	}

	utils.SendEmbedInteractionMsg(s, i, embed)
}
