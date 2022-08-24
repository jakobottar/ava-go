package main

import (
	"ava-go/handlers"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func start() (err error) {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		return errors.New("couldn't find environment variable $BOT_TOKEN")
	}

	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		return errors.New("bot auth: " + err.Error())
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
		return errors.New("bot open: " + err.Error())
	}

	log.Println("bot is running!")

	// add commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(handlers.Commands))
	for idx, cmd := range handlers.Commands {
		regCmd, err := goBot.ApplicationCommandCreate(goBot.State.User.ID, "379276406326165515", cmd)
		if err != nil {
			return errors.New("cannot create command " + cmd.Name + ": " + err.Error())
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
		if err := goBot.ApplicationCommandDelete(goBot.State.User.ID, "379276406326165515", cmd.ID); err != nil {
			return errors.New("cannot delete command " + cmd.Name + ": " + err.Error())
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
