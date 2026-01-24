package commands

import (
	"HimenoSena"
	"HimenoSena/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms-db/discordbot"
)

func GetAllLevel(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, c *HimenoSena.Config) {
	memberData, err := discordbot.GetServerTopMembersByExp(db, c.MainGuildID, 10)
	if err != nil {
		logrus.Error(err)
		return
	}

	var resultStr strings.Builder
	for i, v := range memberData {
		fmt.Fprintf(&resultStr, "%d. **%s** %d等 加入時間: %s\n", i, v.UserName, v.Level, v.JoinAt.Format("2006-01-02"))
	}

	utils.SendInteractionMsg(s, i, resultStr.String())
}
