package bot

import (
	"ava-go/config"
	"ava-go/handlers"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var BotId string

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.MessageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}
