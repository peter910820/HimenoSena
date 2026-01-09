package commands

import (
	"HimenoSena/utils"
	"bufio"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func GetLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const lastNLines = 20
	file, err := os.Open(os.Getenv("KUROHELPER_LOG_PATH"))
	if err != nil {
		logrus.Error(err)
		utils.SendInteractionMsg(s, i, "")
		return
	}
	defer file.Close()

	const maxLineLength = 250
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// 如果一行超過250字元，只取前250字元並加上...
		if len(line) > maxLineLength {
			line = line[:maxLineLength] + "..."
		}
		lines = append(lines, line)
		if len(lines) > lastNLines {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Error(err)
		utils.SendInteractionMsg(s, i, "")
		return
	}

	const linesPerField = 5
	var fields []*discordgo.MessageEmbedField

	// 每5行一個 field
	for i := 0; i < len(lines); i += linesPerField {
		end := i + linesPerField
		if end > len(lines) {
			end = len(lines)
		}

		fieldLines := lines[i:end]
		fieldValue := strings.Join(fieldLines, "\n")

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "",
			Value:  fieldValue,
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:  "KuroHelper Log紀錄(最後10筆)",
		Color:  0x60373E,
		Fields: fields,
	}

	utils.SendEmbedInteractionMsg(s, i, embed)
}
