package utils

import (
	"strings"

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

func SendInteractionMsg(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	if strings.TrimSpace(msg) == "" {
		msg = "該功能目前異常，請稍後再試"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
