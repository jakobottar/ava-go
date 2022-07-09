package bot

import (
	"ava-go/handlers"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Start() {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalln("\u001b[31mERROR:\u001b[0m Couldn't find environment variable $BOT_TOKEN")
		return
	}

	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.MessageHandler)
	goBot.AddHandler(handlers.VoiceStateHandler)

	err = goBot.Open()
	if err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}

	log.Println("bot is running!")
}
