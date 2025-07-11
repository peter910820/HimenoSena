package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena/commands"
	"HimenoSena/event"
	"HimenoSena/models"
	"HimenoSena/utils"
)

var (
	// management database connect
	dbs           = make(map[string]*gorm.DB)
	err           error
	serverUserExp models.ServerUserExp = models.ServerUserExp{}
)

func init() {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// load .env
	err = godotenv.Load(".env")
	if err != nil {
		logrus.Fatal(err)
	}
	// init database
	dbName, db := models.InitDsn()
	dbs[dbName] = db
	models.Migration(dbName, dbs[dbName])
}

func main() {
	c := models.Config{
		Token:         os.Getenv("TOKEN"),
		AppID:         os.Getenv("APP_ID"),
		MainGuildID:   os.Getenv("MAIN_GUILD_ID"),
		BotChannelID:  os.Getenv("BOT_CHANNEL_ID"),
		BotChannelID2: os.Getenv("BOT_CHANNEL_ID2"),
		VoiceManageID: os.Getenv("VOICE_MANAGE_ID"),
		DevCategoryID: os.Getenv("DEV_CATEGORY_ID"),
	}

	c.Bot, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		logrus.Fatal(err)
	}
	c.Bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	c.Bot.AddHandler(ready)
	c.Bot.AddHandler(onInteraction)
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		event.MessageHandler(s, m, &c, dbs[os.Getenv("DATABASE_NAME")], &serverUserExp)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		event.VoiceHandler(s, v, &c)
	})

	err = c.Bot.Open()
	if err != nil {
		logrus.Fatal(err)
	}

	utils.SetUserData(&c, dbs[os.Getenv("DATABASE_NAME")])
	utils.GenerateServerUserExp(&c, dbs[os.Getenv("DATABASE_NAME")], &serverUserExp)

	logrus.Info("bot is now running. Press CTRL+C to exit.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	interruptSignal := <-ch
	c.Bot.Close()
	logrus.Info(interruptSignal)
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "想要傳達給你的愛戀")
	commands.BasicCommand(s)
}

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		delay := s.HeartbeatLatency()
		go commands.Ping(s, i, delay)
	case "send":
		amount := i.ApplicationCommandData().Options[0].StringValue()
		go commands.Send(s, i, amount)
	case "取得身分組":
		go commands.GetRoles(s, i)
	default:
		logrus.Warn("command not founds")
	}
}
