package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"HimenoSena/internal/commands"
)

type Config struct {
	Bot   *discordgo.Session
	Token string
	AppId string
}

func main() {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// load .env
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal(err)
	}

	c := Config{
		Token: os.Getenv("TOKEN"),
		AppId: os.Getenv("Application_ID"),
	}

	c.Bot, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		logrus.Fatal(err)
	}

	c.Bot.AddHandler(ready)
	c.Bot.AddHandler(onInteraction)

	err = c.Bot.Open()
	if err != nil {
		logrus.Fatal(err)
	}

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
	default:
		logrus.Warn("command not founds")
	}
}
