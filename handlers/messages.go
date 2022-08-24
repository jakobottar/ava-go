// Functions handling messages and message-related things
package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// /ping driver function, responds to verify bot life
func ping(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// respond to interaction with success message
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong!",
		},
	})
}

// /glizzy command driver function, prints glizzyL, `content`, glizzyR
func glizzy(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// get map of options
	optionMap := mapOptions(interaction)

	msg := "<a:glizzyR:991176701063221338>" + optionMap["content"].StringValue() + "<a:glizzyL:991176582402150531>"
	if _, err := session.ChannelMessageSend(interaction.ChannelID, msg); err != nil {
		log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}

	// respond to interaction with success message
	//* I don't think I can not respond, but at least I can hide the response
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":hotdog:",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// echo command driver function, echos back `content`
func echo(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// get map of options
	optionMap := mapOptions(interaction)

	// return a message back to the channel, echoing whatever is in the args
	if _, err := session.ChannelMessageSend(interaction.ChannelID, optionMap["content"].StringValue()); err != nil {
		log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}

	log.Println("echo: echoing '", optionMap["content"].StringValue(), "'")

	// respond to interaction with success message
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Echoed!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// remindme command driver function
func remindMe(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	// get map of options
	optionMap := mapOptions(interaction)

	// set timer length
	timerLength := optionMap["time"].IntValue()

	// convert unit arg to time.Duration value
	var unit time.Duration
	switch optionMap["unit"].StringValue() {
	case "min":
		unit = time.Minute
	case "sec":
		unit = time.Second
	case "hr":
		unit = time.Hour
	default:
		return
	}

	// set reminder mention
	var mention *discordgo.User
	if interaction.Member.User != nil {
		// member is only populated on servers
		mention = interaction.Member.User
	} else {
		// user is populated in DMs
		mention = interaction.User
	}

	pronoun := "you"
	if optionMap["user"] != nil {
		mention = optionMap["user"].UserValue(session)
		pronoun = "them"
	}

	// set reminder message
	reminderMessage := ":wave:"
	if optionMap["message"] != nil {
		reminderMessage = optionMap["message"].StringValue()
	}

	// respond to interaction with success message
	responseMessage := fmt.Sprintf("Okay, I will remind %s in %d %s.", pronoun, timerLength, optionMap["unit"].StringValue())
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})

	time.Sleep(time.Duration(timerLength) * unit) // wait

	// set reminder message
	if _, err := session.ChannelMessageSend(interaction.ChannelID, mention.Mention()+" "+reminderMessage); err != nil {
		log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}
	log.Printf("remindme: reminder for %s sent\n", mention.Username)
}
