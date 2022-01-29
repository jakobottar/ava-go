package bot

import (
	"ava-go/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if the message is from a bot, return (prevent loops)
	if m.Author.Bot {
		return
	}

	// ping command - check for life
	if m.Content == "ping" {
		// mention the sender back and say pong!
		_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+", pong!")
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// if the first character of the message matches our bot prefix
	if m.Content[0:1] == config.BotPrefix {
		words := strings.Fields(m.Content) // split message string on space
		command := words[0][1:]            // get first word and remove BotPrefix
		args := words[1:]                  // get the rest of the message

		switch command { // switch on command word
		case "echo": // echo the message that was sent
			echo(s, m, args)

		default: // if the command does not match an existing one, return
			return
		}
	}
}

func echo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// return a message back to the channel, echoing whatever is in the args
	_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
	if err != nil {
		fmt.Println(err.Error())
	}
}
