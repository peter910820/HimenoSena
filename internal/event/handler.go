package event

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func VoiceHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	logrus.Debugf("%v", v)
}
