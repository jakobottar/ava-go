package bot

import (
	"ava-go/config"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

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
			return
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

		case "remindme", "remind":
			remindMe(s, m, args)

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
		return
	}
}

// remindme command driver function, takes commands in the order: !remindme <username> <duration> <unit> <message>
func remindMe(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	//TODO: handle out-of-order args
	// "remind @someone about something in 15 min"
	// "remind 15 min to go get the recycling"

	//TODO: display some error message when parsing fails

	var toRemind *discordgo.User
	var startArgs int = 0

	// if the first arg is a user @,
	id := string(regexp.MustCompile(`\d{18}`).Find([]byte(args[0])))
	if id != "" { // get their user struct from their id
		var err error
		toRemind, err = s.User(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		startArgs++
	} else {
		toRemind = m.Author
	}

	// convert time arg to integer value
	timerLength, err := strconv.Atoi(args[startArgs])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// convert duration arg to time.Duration value
	var unit time.Duration
	var unitStr string

	switch strings.ToLower(args[startArgs+1]) {
	case "minute", "min", "m":
		unit = time.Minute
		unitStr = "minutes"
	case "second", "sec", "s":
		unit = time.Second
		unitStr = "seconds"
	case "hour", "hr", "h":
		unit = time.Hour
		unitStr = "hours"
	default:
		return
	}

	// everything else is the message
	var message string = ":wave:"
	if len(args) > startArgs+2 {
		message = strings.Join(args[startArgs+2:], " ")
	}

	// send confirmation message
	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Reminder set for %d %s.", timerLength, unitStr))
	time.Sleep(time.Duration(timerLength) * unit) // wait
	// set reminder message
	_, _ = s.ChannelMessageSend(m.ChannelID, toRemind.Mention()+" "+message)
}
