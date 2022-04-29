package bot

import (
	"ava-go/config"
	"ava-go/handlers"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var BotId string

func Start() {
	f, err := os.OpenFile("bot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.MessageHandler)
	goBot.AddHandler(handlers.VoiceStateHandler)

	err = goBot.Open()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	log.Println("Bot is running!")
}
