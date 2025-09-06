package commands

import (
	"HimenoSena/models"
	"HimenoSena/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetAllLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, c *models.Config) {
	var memberData []models.Member
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
