package commands

import (
	"HimenoSena/utils"
	"bufio"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const lastNLines = 10
	file, err := os.Open(os.Getenv("KUROHELPER_LOG_PATH"))
	if err != nil {
		utils.SendInteractionMsg(s, i, "")
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > lastNLines {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		utils.SendInteractionMsg(s, i, "")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "KuroHelper Log紀錄(最後10筆)",
		Color: 0x60373E,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Log",
				Value:  strings.Join(lines, "\n"),
				Inline: false,
			},
		},
	}

	utils.SendEmbedInteractionMsg(s, i, embed)
}
