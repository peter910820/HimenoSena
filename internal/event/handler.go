package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"HimenoSena/internal/models"
	"HimenoSena/internal/utils"
)

func VoiceHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate, c *models.Config) {
	// for member := range *vs{
	// 	*vs = append(*vs, *v)
	// }
	name, err := utils.GetName(s, v)
	if err != nil {
		logrus.Error(err)
	}
	// a member into a voice channel
	if v.BeforeUpdate == nil {
		s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***加入了***%s***頻道!", name.Username, name.ChannelName))
	} else {
		if v.ChannelID == "" {
			s.ChannelMessageSend(c.VoiceManageID, fmt.Sprintf("***%s***退出了頻道!", name.Username))
		} else {
			// TODO: add other(Deaf, Mute) status hint
		}
	}
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *models.Config) {
	// avoid responding to itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// judge guild id
	if m.GuildID == c.MainGuildID {
		// judge if message is for bot and ChannelID not command channel
		if m.Author.Bot && (m.ChannelID != c.BotChannelID) && len(m.Embeds) > 0 {
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

		}
	}

}
