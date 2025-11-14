package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena/bot"
	"HimenoSena/handlers"
	"HimenoSena/models"
)

var (
	// management database connect
	dbs             = make(map[string]*gorm.DB)
	err             error
	serverMemberExp models.ServerMemberExp = models.ServerMemberExp{}
	c               models.Config
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
	dbName, db := bot.InitDsn()
	dbs[dbName] = db
	bot.Migration(dbName, dbs[dbName])

	parts := strings.Split(os.Getenv("MAIN_GUILD_ID"), ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	bot.IDMap = make(map[string]struct{}, len(parts))
	for _, id := range parts {
		bot.IDMap[id] = struct{}{}
	}

	c = models.Config{
		Token:          os.Getenv("TOKEN"),
		AppID:          os.Getenv("APP_ID"),
		MainGuildID:    parts, // Not used
		BotChannelID:   os.Getenv("BOT_CHANNEL_ID"),
		BotChannelID2:  os.Getenv("BOT_CHANNEL_ID2"),
		LevelUpChannel: os.Getenv("LEVELUP_CHANNEL"),
		VoiceManageID:  os.Getenv("VOICE_MANAGE_ID"),
		DevCategoryID:  os.Getenv("DEV_CATEGORY_ID"),
	}
}

func main() {
	err = bot.RestoreJsonData(c.MainGuildID[0], &serverMemberExp)
	if err != nil {
		logrus.Fatal(err)
	}

	c.Bot, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		logrus.Fatal(err)
	}
	c.Bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuildVoiceStates

	c.Bot.AddHandler(handlers.Ready)
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		handlers.OnInteractionHandler(s, m, dbs[os.Getenv("DATABASE_NAME")], &serverMemberExp, &c)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.MessageEventHandler(s, m, &c, dbs[os.Getenv("DATABASE_NAME")], &serverMemberExp)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		handlers.VoiceEventHandler(s, v, &c)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		handlers.GuildMemberAddEventHandler(s, m, &c, dbs[os.Getenv("DATABASE_NAME")], &serverMemberExp)
	})

	err = c.Bot.Open()
	if err != nil {
		logrus.Fatal(err)
	}

	bot.SetUserData(&c, dbs[os.Getenv("DATABASE_NAME")])
	bot.GenerateServerUserExp(&c, dbs[os.Getenv("DATABASE_NAME")], &serverMemberExp)

	logrus.Info("bot is now running. Press CTRL+C to exit.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	interruptSignal := <-ch
	c.Bot.Close()
	bot.SaveMemberData(&serverMemberExp)
	logrus.Info(interruptSignal)
}
