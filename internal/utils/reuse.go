package utils

import (
	"github.com/bwmarrin/discordgo"
)

func GetUserName(s *discordgo.Session, v *discordgo.VoiceStateUpdate) (*string, error) {
	name := ""
	su, err := s.User(v.UserID)
	if err != nil {
		return &name, err
	}
	name = su.Username

	return &name, nil
}

func GetChannelName(s *discordgo.Session, v *discordgo.VoiceStateUpdate) (*string, error) {
	name := ""
	sc, err := s.Channel(v.ChannelID)
	if err != nil {
		return &name, err
	}
	name = sc.Name

	return &name, nil
}
