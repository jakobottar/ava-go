package bot

import (
	"ava-go/config"
	"ava-go/handlers"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var BotId string

func Start() (goBot *discordgo.Session, err error) {
	goBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create bot")
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.MessageHandler)
	goBot.AddHandler(handlers.VoiceStateHandler)

	if err = goBot.Open(); err != nil {
		return nil, errors.Wrapf(err, "Failed to open bot")
	}

	log.Println("bot is running!")

	return
}
