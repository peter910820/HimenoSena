package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena"
	"HimenoSena/bot"
	"HimenoSena/utils"

	discordbotdb "github.com/peter910820/discordbot-db"
)

var updateMsgIdTmp = make(map[string]struct{})

func MessageEventHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *HimenoSena.Config, db *gorm.DB, serverUserExp *HimenoSena.ServerMemberExp) {
	// avoid responding to itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		logrus.Error(err)
		return
	}

	// judge guild id
	if m.GuildID == c.MainGuildID {
		// judge if message is for bot and ChannelID not command channel
		if m.Author.Bot && (m.ChannelID != c.BotChannelID && m.ChannelID != c.BotChannelID2) && len(m.Embeds) > 0 && channel.ParentID != c.DevCategoryID {
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

		} else if m.Author.Bot && (m.ChannelID != c.BotChannelID && m.ChannelID != c.BotChannelID2) && channel.ParentID != c.DevCategoryID {
			if !utils.IsDeferredInteraction(m.Message) {
				_, err = s.ChannelMessageSend(c.BotChannelID, fmt.Sprintf("轉送***%s***的訊息:\n%s", m.Author.Username, m.Content))
				if err != nil {
					logrus.Error(err)
				}

				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
		// handle exp feature
		serverUserExp.Mu.Lock()
		defer serverUserExp.Mu.Unlock()
		val, ok := serverUserExp.MemberData[m.Author.ID]
		if ok {
			if val-1 == 0 {
				levelUpExp, level, err := bot.ModifyArticle(m.Author.ID, db)
				if err != nil {
					logrus.Error(err)
					return
				}
				serverUserExp.MemberData[m.Author.ID] = levelUpExp
				embed := &discordgo.MessageEmbed{
					Title:       "🎉**有人升級了**",
					Color:       0xFC9F4D,
					Description: fmt.Sprintf("**%s** 聊天等級升到 **%d** 等", m.Author.GlobalName, level),
				}
				_, err = s.ChannelMessageSendEmbed(c.LevelUpChannel, embed)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			serverUserExp.MemberData[m.Author.ID] = val - 1
		}
	}
}

func VoiceEventHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate, c *HimenoSena.Config) {
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
		return
	}

	// a member switch voice channel
	if v.BeforeUpdate.ChannelID != v.ChannelID {
		channelName, err := utils.GetChannelName(s, v)
		if err != nil {
			logrus.Error(err)
		}
		if strings.TrimSpace(*channelName) == "" {
			s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***退出語音頻道!", *username))
		} else {
			s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***跑到***%s***頻道了!", *username, *channelName))
		}
		return
	}

	//  a member leave a voice channel
	if v.ChannelID == "" {
		s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***退出了頻道!", *username))
		return
	}

	// deaf event
	if v.BeforeUpdate.SelfDeaf != v.SelfDeaf {
		s.ChannelMessageSend(c.VoiceManageID, func() string {
			if v.SelfDeaf {
				return fmt.Sprintf("***%s***拒聽中!", *username)
			}
			return fmt.Sprintf("***%s***解除拒聽!", *username)
		}())
		return
	}

	// mute event
	if v.BeforeUpdate.SelfMute != v.SelfMute {
		s.ChannelMessageSend(c.VoiceManageID, func() string {
			if v.SelfMute {
				return fmt.Sprintf("***%s***靜音中!", *username)
			}
			return fmt.Sprintf("***%s***解除靜音!", *username)
		}())
		return
	}

	// steam event
	if v.BeforeUpdate.SelfStream != v.SelfStream {
		s.ChannelMessageSend(c.VoiceManageID, func() string {
			if v.SelfStream {
				return fmt.Sprintf("***%s***開始直播!", *username)
			}
			return fmt.Sprintf("***%s***關掉直播了QQ!", *username)
		}())
		return
	}

	// other event
	s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***其他語音狀態改變!", *username))
}

func GuildMemberAddEventHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd, c *HimenoSena.Config, db *gorm.DB, serverUserExp *HimenoSena.ServerMemberExp) {
	serverUserExp.Mu.Lock()
	defer serverUserExp.Mu.Unlock()
	err := discordbotdb.CreateMember(db, discordbotdb.Member{
		UserID:   m.User.ID,
		ServerID: c.MainGuildID,
		UserName: m.User.Username,
		JoinAt:   m.JoinedAt,
	})
	if err != nil {
		logrus.Error(err)
	}

	_, ok := serverUserExp.MemberData[m.Member.User.ID]
	if !ok {
		serverUserExp.MemberData[m.Member.User.ID] = 5
	}
}

func MessageUpdateHandler(s *discordgo.Session, m *discordgo.MessageUpdate, c *HimenoSena.Config) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 只轉送自己群組的訊息
	if m.GuildID == c.MainGuildID {
		if m.Author.Bot && (m.ChannelID != c.BotChannelID && m.ChannelID != c.BotChannelID2) && channel.ParentID != c.DevCategoryID {
			if _, ok := updateMsgIdTmp[m.ID]; ok {
				delete(updateMsgIdTmp, m.ID)
				return
			}
			embeds := make([]*discordgo.MessageEmbed, 0, len(m.Embeds))
			for _, e := range m.Embeds {
				embeds = append(embeds, &discordgo.MessageEmbed{
					Title:       e.Title,
					Description: e.Description,
					URL:         e.URL,
					Timestamp:   e.Timestamp,
					Color:       e.Color,
					Footer:      e.Footer,
					Image:       e.Image,
					Thumbnail:   e.Thumbnail,
					Author:      e.Author,
					Fields:      e.Fields,
				})
			}

			returnMsg := &discordgo.MessageSend{
				Content: fmt.Sprintf(
					"轉送 ***%s*** 的訊息：\n%s",
					m.Author.Username,
					m.Content,
				),

				Embeds: embeds,
			}

			if _, err := s.ChannelMessageSendComplex(c.BotChannelID, returnMsg); err != nil {
				logrus.Error(err)
			}

			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				logrus.Error(err)
			}

			updateMsgIdTmp[m.Message.ID] = struct{}{}
		}
	}
}
