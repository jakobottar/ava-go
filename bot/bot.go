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

// const (
// 	membersPageLimit = 1000
// )

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// declare intents (needed to be able to get member info)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

// with help from https://github.com/ewohltman/ephemeral-roles/blob/6bbd1e38824b73df892a269657f967eeab583e46/internal/pkg/operations/operations.go#L258
//* I don't think I need any of this, it appears to be redundant? Or when session.State.Guild(id) fails?
// func updateStateGuilds(session *discordgo.Session, guildID string) (*discordgo.Guild, error) {
// 	guild, err := session.Guild(guildID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	roles, err := session.GuildRoles(guildID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	channels, err := session.GuildChannels(guildID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	members, err := getGuildMembers(session, guildID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	guild.Roles = roles
// 	guild.Channels = channels
// 	guild.Members = members
// 	guild.MemberCount = len(members)

// 	err = session.State.GuildAdd(guild)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return guild, nil
// }

// func getGuildMembers(session *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
// 	var members []*discordgo.Member // initialize empty array to hold all members
// 	for {
// 		membersPage, err := session.GuildMembers(guildID, "", membersPageLimit) // get a page of members
// 		if err != nil {
// 			return nil, err
// 		}

// 		members = append(members, membersPage...) // append the page to the array of members

// 		if len(membersPage) < membersPageLimit { // if we got less than a full page of members, exit
// 			break
// 		}
// 	}

// 	return members, nil
// }

// func getGuild(session *discordgo.Session, guildID string) (*discordgo.Guild, error) {
// 	guild, err := session.State.Guild(guildID)
// 	if err != nil {
// 		guild, err := updateStateGuilds(session, guildID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return guild, nil
// 	}

// 	return guild, nil
// }

func messageHandler(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// if the message is from a bot, return (prevent loops)
	if msg.Author.Bot {
		return
	}

	// ping command - check for life
	if msg.Content == "ping" {
		// mention the sender back and say pong!
		_, err := session.ChannelMessageSend(msg.ChannelID, msg.Author.Mention()+", pong!")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		guild, err := session.State.Guild(msg.GuildID)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("This server has %d members!\n", guild.MemberCount)
		for i := 0; i < len(guild.Members); i++ {
			fmt.Println(guild.Members[i].User)
		}
	}

	// if the first character of the message matches our bot prefix
	if msg.Content[0:1] == config.BotPrefix {
		words := strings.Fields(msg.Content) // split message string on space
		command := words[0][1:]              // get first word and remove BotPrefix
		args := words[1:]                    // get the rest of the message

		switch command { // switch on command word
		case "echo": // echo the message that was sent
			echo(session, msg, args)

		case "remindme", "remind":
			remindMe(session, msg, args)

		default: // if the command does not match an existing one, return
			return
		}
	}
}

func echo(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	// return a message back to the channel, echoing whatever is in the args
	_, err := session.ChannelMessageSend(msg.ChannelID, strings.Join(args, " "))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// remindme command driver function, takes commands in the order: !remindme <username> <duration> <unit> <message>
func remindMe(session *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	//TODO: handle out-of-order args
	// "remind @someone about something in 15 min"
	// "remind 15 min to go get the recycling"

	//TODO: display some error message when parsing fails

	var toRemind *discordgo.User
	var startArgs int = 0

	// if the first arg is a user @,
	id := string(regexp.MustCompile(`\d{18}`).Find([]byte(args[0]))) //! if args is empty, we fault here
	if id != "" {                                                    // get their user struct from their id
		var err error
		toRemind, err = session.User(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		startArgs++
	} else {
		toRemind = msg.Author
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
	_, _ = session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Reminder set for %d %s.", timerLength, unitStr))
	time.Sleep(time.Duration(timerLength) * unit) // wait
	// set reminder message
	_, _ = session.ChannelMessageSend(msg.ChannelID, toRemind.Mention()+" "+message)
}
