package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"HimenoSena/internal/models"
)

func VoiceHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate, c *models.Config) {
	if v.BeforeUpdate == nil {

	}
	logrus.Debugf("%s", v.Member.User.Username)

	// s.ChannelMessageSend()
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
