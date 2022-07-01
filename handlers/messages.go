package handlers

import (
	"ava-go/config"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// if the message is from a bot, return (prevent loops)
	//  or if the message is only an image, it's content is empty
	if msg.Author.Bot || msg.Content == "" {
		return
	}

	// ping command - check for life
	if msg.Content == "ping" {
		log.Println("msghandler: caught ping command")
		// mention the sender back and say pong!
		_, err := session.ChannelMessageSend(msg.ChannelID, msg.Author.Mention()+", pong!")
		if err != nil {
			log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			return
		}
	}

	// if the first character of the message matches our bot prefix
	if msg.Content[0:1] == config.BotPrefix {
		words := strings.Fields(msg.Content) // split message string on space
		command := words[0][1:]              // get first word and remove BotPrefix
		args := words[1:]                    // get the rest of the message

		switch command { // switch on command word
		case "echo": // echo the message that was sent
			log.Println("msghandler: caught echo command")
			echo(session, msg, args)
		
		case "glizzy": // glizzy
			log.Println("msghandler: caught glizzy command")
			// TODO: get these emojis not in a hard-coded way
			_, err := session.ChannelMessageSend(msg.ChannelID, "<a:glizzyR:991176701063221338>" + strings.Join(args, " ") + "<a:glizzyL:991176582402150531>" )
			if err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
				return
			}

		case "remindme", "remind":
			log.Println("msghandler: caught remindme command")
			remindMe(session, msg, args)

		case "shuffle":
			log.Println("msghandler: caught shuffle command")
			shuffleVCs(session, msg)

		default: // if the command does not match an existing one, return
			return
		}
	}
}

func echo(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// return a message back to the channel, echoing whatever is in the args
	_, err := session.ChannelMessageSend(msg.ChannelID, strings.Join(args, " "))
	if err != nil {
		log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
		return
	}

	log.Println("echo: echoing '", strings.Join(args, " "), "'")
}

// remindme command driver function, takes commands in the order: !remindme <username> <duration> <unit> <message>
func remindMe(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		return
	}

	log.Println("remindme: parsing 'remindme' args...")

	var timerLength int = 1                  // timer length
	var mention *discordgo.User = msg.Author // find user to mention
	var unitStr string = "m"                 // find timer unit
	var message string = ":wave:"            // find message

	var mentionName []string
	var err error
	var unitIdx int

	for i, arg := range args {
		timerLength, err = strconv.Atoi(arg)
		if err == nil {
			unitIdx = i + 1
			break
		}
		mentionName = append(mentionName, arg)
	}

	if len(mentionName) == 0 { // if there's nothing before the duration
		mention = msg.Author
	} else { // find the mentioned user
		name := strings.Join(mentionName, " ")
		id := string(regexp.MustCompile(`\d{18}`).Find([]byte(name)))
		if id != "" { // if 'name' is a user id or mention
			mention, err = session.User(id)
			if err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
				return
			}
		} else { // if 'name' is nickname or username
			guild, err := session.State.Guild(msg.GuildID)
			if err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
				return
			}

			for _, member := range guild.Members {
				if name == member.User.Username || name == member.Nick {
					mention = member.User
					break
				}
			}
		}
	}

	// convert duration arg to time.Duration value
	unitStr = strings.ToLower(args[unitIdx])
	var unit time.Duration
	switch unitStr {
	case "minutes", "minute", "mins", "min", "m":
		unit = time.Minute
		unitStr = "minutes"
	case "seconds", "second", "secs", "sec", "s":
		unit = time.Second
		unitStr = "seconds"
	case "hours", "hour", "hrs", "hr", "h":
		unit = time.Hour
		unitStr = "hours"
	default:
		return
	}

	if unitIdx == len(args)-1 {
		message = ":wave:"
	} else {
		message = strings.Join(args[unitIdx+1:], " ")
	}

	log.Printf("remindme: set for %s in %d %s\n", mention.Username, timerLength, unitStr)

	// send confirmation message
	_, _ = session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Reminder set for %d %s.", timerLength, unitStr))
	time.Sleep(time.Duration(timerLength) * unit) // wait
	// set reminder message
	_, _ = session.ChannelMessageSend(msg.ChannelID, mention.Mention()+" "+message)
	log.Printf("remindme: reminder for %s sent\n", mention.Username)
}
