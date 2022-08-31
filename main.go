package main

import (
	"ava-go/handlers"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func start() error {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		return errors.New("start: couldn't find environment variable $BOT_TOKEN")
	}

	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("start: bot create %s", err.Error())
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(handlers.VoiceStateHandler)

	// add slash command handlers
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	if err = goBot.Open(); err != nil {
		return fmt.Errorf("start: bot open %s", err.Error())
	}

	log.Println("bot is running!")

	// lookup guild id from .env
	guildID, ok := os.LookupEnv("GUILD_ID")
	if !ok {
		return errors.New("start: couldn't find environment variable $GUILD_ID")
	}

	// add commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(handlers.Commands))

	for idx, cmd := range handlers.Commands {
		if cmd.Name == "glizzy" {
			// fetch glizzy emojis
			if err := handlers.FetchGlizzy(goBot, guildID); err != nil {
				log.Println("\033[33mWARNING:\033[0m skipping glizzy -", err.Error())
				continue
			}
		}

		regCmd, err := goBot.ApplicationCommandCreate(goBot.State.User.ID, guildID, cmd)
		if err != nil {
			return fmt.Errorf("start: cannot create command %s: %s", cmd.Name, err.Error())
		}

		registeredCommands[idx] = regCmd
	}

	defer goBot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)    // stop on ctrl-C
	signal.Notify(stop, syscall.SIGTERM) // stop on docker stop
	<-stop

	// clean up added commands
	for _, cmd := range registeredCommands {
		if cmd != nil { // in case we skip something when registering commands
			if err := goBot.ApplicationCommandDelete(goBot.State.User.ID, guildID, cmd.ID); err != nil {
				return fmt.Errorf("start: cannot delete command %s: %s", cmd.Name, err.Error())
			}
		}
	}

	log.Println("shutting down gracefully!")

	return nil
}

func main() {
	if err := start(); err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
	}
}
