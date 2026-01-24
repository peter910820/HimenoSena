package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"HimenoSena"
	"HimenoSena/bot"
	"HimenoSena/handlers"

	seaottermsdb "seaotterms-db"
)

var (
	// management database connect
	dbs             *gorm.DB
	err             error
	serverMemberExp HimenoSena.ServerMemberExp = HimenoSena.ServerMemberExp{}
	c               HimenoSena.Config
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
	portStr := os.Getenv("DATABASE_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logrus.Fatal(err)
	}

	// 初始化DB
	db, err := seaottermsdb.InitDsn(seaottermsdb.ConnectDBConfig{
		Owner:    os.Getenv("DATABASE_OWNER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_NAME"),
		Port:     port,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	dbs = db.DB

	// Migration
	seaottermsdb.Migration(db)

	c = HimenoSena.Config{
		Token:          os.Getenv("TOKEN"),
		AppID:          os.Getenv("APP_ID"),
		MainGuildID:    os.Getenv("MAIN_GUILD_ID"),
		BotChannelID:   os.Getenv("BOT_CHANNEL_ID"),
		BotChannelID2:  os.Getenv("BOT_CHANNEL_ID2"),
		LevelUpChannel: os.Getenv("LEVELUP_CHANNEL"),
		VoiceManageID:  os.Getenv("VOICE_MANAGE_ID"),
		DevCategoryID:  os.Getenv("DEV_CATEGORY_ID"),
	}
}

func main() {
	err = bot.RestoreJsonData(c.MainGuildID, &serverMemberExp)
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
		handlers.OnInteractionHandler(s, m, dbs, &serverMemberExp, &c)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.MessageEventHandler(s, m, &c, dbs, &serverMemberExp)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		handlers.VoiceEventHandler(s, v, &c)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		handlers.GuildMemberAddEventHandler(s, m, &c, dbs, &serverMemberExp)
	})
	c.Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		handlers.MessageUpdateHandler(s, m, &c)
	})

	err = c.Bot.Open()
	if err != nil {
		logrus.Fatal(err)
	}

	bot.SetUserData(&c, dbs)
	bot.GenerateServerUserExp(&c, dbs, &serverMemberExp)

	logrus.Info("bot is now running. Press CTRL+C to exit.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	interruptSignal := <-ch
	c.Bot.Close()
	bot.SaveMemberData(&serverMemberExp)
	logrus.Info(interruptSignal)
}
