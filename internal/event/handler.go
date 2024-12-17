package event

import (
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
