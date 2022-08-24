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
	}

	CommandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"ping": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught ping command")
			ping(session, interaction)
		},
		"glizzy": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught glizzy command")
			glizzy(session, interaction)

		},
		"shuffle": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught shuffle command")
			shuffle(session, interaction)
		},
	}
)

// convert slashcommand options to a map object
func mapOptions(interaction *discordgo.InteractionCreate) (optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	// get map of options
	options := interaction.ApplicationCommandData().Options
	optionMap = make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return
}
