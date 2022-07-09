package bot

import (
	"ava-go/handlers"
	"errors"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Start() (err error) {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		return errors.New("couldn't find environment variable $BOT_TOKEN")
	}

	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		return
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.MessageHandler)
	goBot.AddHandler(handlers.VoiceStateHandler)

	if err = goBot.Open(); err != nil {
		return
	}

	log.Println("bot is running!")

	return nil
}
