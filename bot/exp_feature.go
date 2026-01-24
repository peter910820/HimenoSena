package bot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena"

	"seaotterms-db/discordbot"
)

// set user data into database
func SetUserData(c *HimenoSena.Config, db *gorm.DB) {
	members, err := c.Bot.GuildMembers(c.MainGuildID, "", 1000)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, member := range members {
		if !member.User.Bot {
			err := discordbot.CreateMember(db, discordbot.Member{
				UserID:   member.User.ID,
				ServerID: c.MainGuildID,
				UserName: member.User.Username,
				JoinAt:   member.JoinedAt,
			})
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func GenerateServerUserExp(c *HimenoSena.Config, db *gorm.DB, serverUserExp *HimenoSena.ServerMemberExp) {
	if len(serverUserExp.MemberData) != 0 {
		members, err := discordbot.QueryMembers(db)
		if err != nil {
			logrus.Fatal(err)
		}

		for _, member := range members {
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
		members, err := discordbot.QueryMembers(db)
		if err != nil {
			logrus.Fatal(err)
		}

		for _, member := range members {
			serverUserExp.MemberData[member.UserID] = member.LevelUpExp
		}
	}
}

func ModifyArticle(userID string, db *gorm.DB) (uint, uint, error) {
	var resultLevelUpExp uint
	var resultLevel uint

	err := db.Transaction(func(tx *gorm.DB) error {

		member, err := discordbot.QueryMemberByUserID(tx, userID)
		if err != nil {
			return err
		}

		levelUpExp := 5 + (member.Level+1)*2 - 2
		data := discordbot.Member{
			Level:      member.Level + 1,
			Exp:        member.Exp + 5 + (member.Level)*2 - 2,
			LevelUpExp: levelUpExp,
			UpdatedAt:  time.Now(),
		}

		if err := discordbot.UpdateMemberLevel(tx, userID, data); err != nil {
			return err
		}

		resultLevelUpExp = levelUpExp
		resultLevel = member.Level + 1

		return nil
	})

	if err != nil {
		logrus.Error(err)
		return 0, 0, err
	}

	return resultLevelUpExp, resultLevel, nil
}

func SaveMemberData(data *HimenoSena.ServerMemberExp) {
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

func RestoreJsonData(mainGuildID string, serverMemberExp *HimenoSena.ServerMemberExp) error {
	path := mainGuildID + "_memberData.json"

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var memberData HimenoSena.ServerMemberExp
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
