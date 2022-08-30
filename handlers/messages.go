// Functions handling messages and message-related things
package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	glizzyR *discordgo.Emoji
	glizzyL *discordgo.Emoji
)

// /ping driver function, responds to verify bot life
func ping(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	// respond to interaction with success message
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong!",
		},
	})

	return err
}

// get glizzy emotes from server for use in the /glizzy command
func FetchGlizzy(session *discordgo.Session, guildID string) error {
	guild, err := session.Guild(guildID)
	if err != nil {
		return fmt.Errorf("fetchglizzy: could not get guild ID %s", err.Error())
	}

	// TODO: catch if we don't find both emotes, throw error
	for _, emoji := range guild.Emojis {
		if emoji.Name == "glizzyR" {
			glizzyR = emoji
		} else if emoji.Name == "glizzyL" {
			glizzyL = emoji
		}
	}

	return nil
}

// /glizzy command driver function, prints glizzyL, `content`, glizzyR
func glizzy(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	// get map of options
	optionMap := mapOptions(interaction)

	msg := fmt.Sprintf("%s%s%s", glizzyR.MessageFormat(), optionMap["content"].StringValue(), glizzyL.MessageFormat())
	if _, err := session.ChannelMessageSend(interaction.ChannelID, msg); err != nil {
		return fmt.Errorf("glizzy: %s", err.Error())
	}

	// respond to interaction with success message
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":hotdog:",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return err
}

// echo command driver function, echos back `content`
func echo(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	// get map of options
	optionMap := mapOptions(interaction)

	// return a message back to the channel, echoing whatever is in the args
	if _, err := session.ChannelMessageSend(interaction.ChannelID, optionMap["content"].StringValue()); err != nil {
		return fmt.Errorf("echo: could not send message %s", err.Error())
	}

	log.Printf("echo: echoing '%s'\n", optionMap["content"].StringValue())

	// respond to interaction with success message
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Echoed!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return err
}

// remindme command driver function
func remindMe(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
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
		return fmt.Errorf("remindMe: Invalid Time Unit '%s'", optionMap["unit"].StringValue())
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
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})

	if err != nil {
		return fmt.Errorf("remindMe: could not send message %s", err.Error())
	}

	time.Sleep(time.Duration(timerLength) * unit) // wait

	// set reminder message
	if _, err := session.ChannelMessageSend(interaction.ChannelID, mention.Mention()+" "+reminderMessage); err != nil {
		return fmt.Errorf("remindMe: could not send message %s", err.Error())
	}

	log.Printf("remindme: reminder for %s sent\n", mention.Username)

	return nil
}
