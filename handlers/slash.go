package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go
var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "check if the bot is running",
		},
		{
			Name:        "glizzy",
			Description: "glizzy gladiators",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "emote or content",
					Required:    true,
				},
			},
		},
		{
			Name:        "shuffle",
			Description: "shuffle voice channels",
		},
		{
			Name:        "echo",
			Description: "echo back `content`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "content to be echoed back",
					Required:    true,
				},
			},
		},
		{
			Name:        "remindme",
			Description: "remind me (or another user) after a given amount of time",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "time",
					Description: "timer length (number)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "unit",
					Description: "time unit",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "seconds",
							Value: "sec",
						},
						{
							Name:  "minutes",
							Value: "min",
						},
						{
							Name:  "hours",
							Value: "hr",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "user",
					Description: "user to be reminded (optional, default @me)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "message to be included in the reminder",
					Required:    false,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught ping command")
			if err := ping(s, i); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		},
		"glizzy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught glizzy command")
			if err := glizzy(s, i); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		},
		"shuffle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught shuffle command")
			if err := shuffle(s, i); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		},
		"echo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught echo command")
			if err := echo(s, i); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		},
		"remindme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught remindme command")
			if err := remindMe(s, i); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		},
	}
)

// convert slashcommand options to a map object
func mapOptions(interaction *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	// get map of options
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}
