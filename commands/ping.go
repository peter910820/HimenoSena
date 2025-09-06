package commands

import (
	"HimenoSena/utils"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate, delay time.Duration) {
	utils.SendInteractionMsg(s, i, fmt.Sprintf("延遲時間為: %v", delay))
}
