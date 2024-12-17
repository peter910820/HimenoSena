package utils

import (
	"github.com/bwmarrin/discordgo"
)

type Name struct {
	Username    string
	ChannelName string
}

func GetName(s *discordgo.Session, v *discordgo.VoiceStateUpdate) (*Name, error) {
	name := Name{}

	su, err := s.User(v.UserID)
	if err != nil {
		return &name, err
	}
	name.Username = su.Username

	sc, err := s.Channel(v.ChannelID)
	if err != nil {
		return &name, err
	}
	name.ChannelName = sc.Name

	return &name, nil
}
