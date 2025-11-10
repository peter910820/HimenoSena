package utils

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	// option not found error
	ErrOptionNotFound = errors.New("option: option not found")
	// option translate fail error
	ErrOptionTranslateFail = errors.New("option: value translate fail")
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

func SendEmbedInteractionMsg(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func GetGuildNick(s *discordgo.Session, guildID string, userID string) (string, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		// 失敗處理
		logrus.Error(err)
		return "", err
	}

	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}
	return nick, nil
}

// get slash command options
func GetOptions(i *discordgo.InteractionCreate, name string) (string, error) {
	for _, v := range i.ApplicationCommandData().Options {
		if v.Name == name {
			value, ok := v.Value.(string) // type assertion
			if !ok {
				return "", ErrOptionTranslateFail
			} else {
				return value, nil
			}
		}
	}
	return "", ErrOptionNotFound
}
