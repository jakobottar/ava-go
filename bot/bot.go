package bot

import (
	"ava-go/config"
	"ava-go/handlers"
	"log"

	"github.com/bwmarrin/discordgo"
)

var BotId string

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
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
