package bot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"HimenoSena/model"
)

// set user data into database
func SetUserData(c *model.Config, db *gorm.DB) {
	members, err := c.Bot.GuildMembers(c.MainGuildID, "", 1000)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, member := range members {
		if !member.User.Bot {
			CreateUser(c, member, db)
		}
	}
}

func CreateUser(c *model.Config, member *discordgo.Member, db *gorm.DB) {
	data := model.Member{
		UserID:   member.User.ID,
		ServerID: c.MainGuildID,
		UserName: member.User.Username,
		JoinAt:   member.JoinedAt,
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&data).Error
	if err != nil {
		logrus.Fatal(err)
	}
}

func GenerateServerUserExp(c *model.Config, db *gorm.DB, serverUserExp *model.ServerMemberExp) {
	if len(serverUserExp.MemberData) != 0 {
		members := queryUser(db)
		for _, member := range *members {
			_, ok := serverUserExp.MemberData[member.UserID]
			if ok {
				continue
			} else {
				serverUserExp.MemberData[member.UserID] = member.LevelUpExp
			}
		}
	} else {
		serverUserExp.ServerID = c.MainGuildID
		serverUserExp.MemberData = make(map[string]uint)
		members := queryUser(db)
		for _, member := range *members {
			serverUserExp.MemberData[member.UserID] = member.LevelUpExp
		}
	}
}

func queryUser(db *gorm.DB) *[]model.Member {
	var UserData []model.Member

	result := db.Find(&UserData)
	if result.Error != nil {
		logrus.Error(result.Error)
	}
	return &UserData
}

// query seingle member fo database use userID
func QueryUser(userID string, db *gorm.DB) (*model.Member, error) {
	var memberData model.Member
	err := db.Select("level, exp").Where("user_id = ?", userID).First(&memberData).Error
	if err != nil {
		return &memberData, err
	}
	return &memberData, nil
}

func ModifyArticle(userID string, db *gorm.DB) (uint, uint, error) {
	var memberData model.Member
	err := db.Select("level, exp").Where("user_id = ?", userID).First(&memberData).Error
	if err != nil {
		logrus.Error(err)
	}
	levelUpExp := 5 + (memberData.Level+1)*2 - 2
	data := model.Member{
		Level:      memberData.Level + 1,
		Exp:        memberData.Exp + 5 + (memberData.Level)*2 - 2,
		LevelUpExp: levelUpExp,
		UpdatedAt:  time.Now(),
	}

	err = db.Model(&model.Member{}).Where("user_id = ?", userID).
		Select("level", "exp", "level_up_exp", "updated_at").Updates(data).Error
	if err != nil {
		return 0, 0, err
	}
	return levelUpExp, memberData.Level + 1, nil
}

func SaveMemberData(data *model.ServerMemberExp) {
	file, err := os.Create(data.ServerID + "_memberData.json")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data.MemberData); err != nil {
		logrus.Error(err)
	}
}

func RestoreJsonData(mainGuildID string, serverMemberExp *model.ServerMemberExp) error {
	path := mainGuildID + "_memberData.json"

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var memberData model.ServerMemberExp
		if err := json.Unmarshal(data, &memberData.MemberData); err != nil {
			return err
		}

		serverMemberExp.ServerID = mainGuildID
		serverMemberExp.MemberData = memberData.MemberData
	} else {
		return nil
	}
	return nil
}
