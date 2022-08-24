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
	}

	CommandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"ping": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught ping command")
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong!",
				},
			})
		},
		"glizzy": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			log.Println("slashhandler: caught glizzy command")

			// get map of options
			optionMap := mapOptions(interaction)

			// TODO: I don't think I can not respond, but at least I can hide the response
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🌭",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			msg := "<a:glizzyR:991176701063221338>" + optionMap["content"].StringValue() + "<a:glizzyL:991176582402150531>"
			if _, err := session.ChannelMessageSend(interaction.ChannelID, msg); err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
				return
			}
		},
	}
)

func mapOptions(interaction *discordgo.InteractionCreate) (optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	// get map of options
	options := interaction.ApplicationCommandData().Options
	optionMap = make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return
}
