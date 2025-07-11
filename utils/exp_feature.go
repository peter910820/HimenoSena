package utils

import (
	"HimenoSena/models"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// set user data into database
func SetUserData(c *models.Config, db *gorm.DB) {
	members, err := c.Bot.GuildMembers(c.MainGuildID, "", 1000)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, member := range members {
		if !member.User.Bot {
			createUser(c, member, db)
		}
	}
}

func createUser(c *models.Config, member *discordgo.Member, db *gorm.DB) {
	data := models.Member{
		MemberID:   member.User.ID,
		ServerID:   c.MainGuildID,
		MemberName: member.User.Username,
		JoinAt:     member.JoinedAt,
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&data).Error
	if err != nil {
		logrus.Fatal(err)
	}
}

func GenerateServerUserExp(c *models.Config, db *gorm.DB, serverUserExp *models.ServerMemberExp) {
	serverUserExp.ServerID = c.MainGuildID
	serverUserExp.MemberData = make(map[string]uint)
	members := queryUser(db)
	for _, member := range *members {
		serverUserExp.MemberData[member.MemberID] = member.LevelUpExp
	}
}

func queryUser(db *gorm.DB) *[]models.Member {
	var UserData []models.Member

	result := db.Find(&UserData)
	if result.Error != nil {
		logrus.Error(result.Error)
	}
	return &UserData
}

func ModifyArticle(memberID string, db *gorm.DB) (uint, uint, error) {
	var memberData models.Member
	err := db.Select("level, exp").Where("user_id = ?", memberID).First(&memberData).Error
	if err != nil {
		logrus.Error(err)
	}
	levelUpExp := 5 + (memberData.Level+1)*2 - 2
	logrus.Debugf("%+v", memberData)
	data := models.Member{
		Level:      memberData.Level + 1,
		Exp:        memberData.Exp + 5 + (memberData.Level)*2 - 2,
		LevelUpExp: levelUpExp,
		UpdatedAt:  time.Now(),
	}

	err = db.Model(&models.Member{}).Where("user_id = ?", memberID).
		Select("level", "exp", "level_up_exp", "updated_at").Updates(data).Error
	if err != nil {
		return 0, 0, err
	}
	return levelUpExp, memberData.Level + 1, nil
}
