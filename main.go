package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Bot   *discordgo.Session
	Token string
	AppId string
}

func main() {
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

	// bot.AddHandler(ready)
	// bot.AddHandler(guildMemberAdd)
	// bot.AddHandler(onInteraction)

	err = c.Bot.Open()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("bot is now running. Press CTRL+C to exit.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	interruptSignal := <-ch
	c.Bot.Close()
	logrus.Println(interruptSignal)
}
