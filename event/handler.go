package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena/models"
	"HimenoSena/utils"
)

func VoiceHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate, c *models.Config) {
	username, err := utils.GetUserName(s, v)
	if err != nil {
		logrus.Error(err)
	}
	// a member into a voice channel
	if v.BeforeUpdate == nil {
		channelName, err := utils.GetChannelName(s, v)
		if err != nil {
			logrus.Error(err)
		}
		s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***加入了***%s***頻道!", *username, *channelName))
	} else {
		if v.ChannelID == "" {
			s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***退出了頻道!", *username))
		} else {
			if v.BeforeUpdate.ChannelID != v.ChannelID {
				channelName, err := utils.GetChannelName(s, v)
				if err != nil {
					logrus.Error(err)
				}
				s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***跑到***%s***頻道了!", *username, *channelName))
			} else if v.BeforeUpdate.SelfDeaf != v.SelfDeaf {
				s.ChannelMessageSend(c.VoiceManageID, func() string {
					if v.SelfDeaf {
						return fmt.Sprintf("***%s***拒聽中!", *username)
					}
					return fmt.Sprintf("***%s***解除拒聽!", *username)
				}())
			} else if v.BeforeUpdate.SelfMute != v.SelfMute {
				s.ChannelMessageSend(c.VoiceManageID, func() string {
					if v.SelfMute {
						return fmt.Sprintf("***%s***靜音中!", *username)
					}
					return fmt.Sprintf("***%s***解除靜音!", *username)
				}())
			} else if v.BeforeUpdate.SelfStream != v.SelfStream {
				s.ChannelMessageSend(c.VoiceManageID, func() string {
					if v.SelfStream {
						return fmt.Sprintf("***%s***開始直播!", *username)
					}
					return fmt.Sprintf("***%s***關掉直播了QQ!", *username)
				}())
			} else {
				s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***語音狀態改變!", *username))
			}
		}
	}
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *models.Config, db *gorm.DB, serverUserExp *models.ServerMemberExp) {
	// avoid responding to itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// channel, err := s.Channel(m.ChannelID)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }
	// // judge DevCategoryID id
	// if channel.ParentID == c.DevCategoryID {
	// 	return
	// }
	// judge guild id
	if m.GuildID == c.MainGuildID {
		// judge if message is for bot and ChannelID not command channel
		if m.Author.Bot && (m.ChannelID != c.BotChannelID && m.ChannelID != c.BotChannelID2) && len(m.Embeds) > 0 {
			// transferMsg := fmt.Sprintf("轉送*%s*的訊息:\n%v", m.Author.Username, m.Content)
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				logrus.Error(err)
			}
			_, err = s.ChannelMessageSend(c.BotChannelID, fmt.Sprintf("轉送***%s***的訊息:", m.Author.Username))
			if err != nil {
				logrus.Error(err)
			}
			_, err = s.ChannelMessageSendEmbed(c.BotChannelID, m.Embeds[0])
			if err != nil {
				logrus.Error(err)
			}

		} else if m.Author.Bot && (m.ChannelID != c.BotChannelID && m.ChannelID != c.BotChannelID2) && len(m.Embeds) == 0 {
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				logrus.Error(err)
			}
			utils.GetGuildNick(s, c.MainGuildID, m.Author.ID)
			_, err = s.ChannelMessageSend(c.BotChannelID, fmt.Sprintf("轉送***%s***的訊息:\n%s", m.Author.Username, m.Content))
			if err != nil {
				logrus.Error(err)
			}
		}
		// handle exp feature
		serverUserExp.Mu.Lock()
		defer serverUserExp.Mu.Unlock()
		val, ok := serverUserExp.MemberData[m.Author.ID]
		if ok {
			if val-1 == 0 {
				levelUpExp, level, err := utils.ModifyArticle(m.Author.ID, db)
				logrus.Debug(levelUpExp)
				if err != nil {
					logrus.Error(err)
					return
				}
				serverUserExp.MemberData[m.Author.ID] = levelUpExp
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** 聊天等級升到 **%d** 等", m.Author.GlobalName, level))
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			serverUserExp.MemberData[m.Author.ID] = val - 1
		}
		logrus.Debug(serverUserExp.MemberData[m.Author.ID])
	}
}
