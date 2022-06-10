package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// TODO: make sure I'm catching (and logging) all errors

const BUFFER_CHANNELS = 2 // number of empty channels to leave

// handler for voice state changes (user join/leave)
func VoiceStateHandler(session *discordgo.Session, voiceState *discordgo.VoiceStateUpdate) {
	guild, _ := session.State.Guild(voiceState.GuildID)
	memberCount := getVCMembers(guild.Channels, guild.VoiceStates)

	// if any channels are empty, count them and remove any extra
	emptyChannels := 0
	for id, numMembers := range memberCount {
		if numMembers == 0 {
			emptyChannels++
			if emptyChannels > BUFFER_CHANNELS {
				_, _ = session.ChannelDelete(id)
				emptyChannels--
			}
		}
	}

	log.Println("voicehandler: cleared extra channels")

	// emptyChannels should only ever be 1-BUFFER_CHANNELS at min so we only need to make one new channel
	if emptyChannels < BUFFER_CHANNELS {
		makeNewVoiceChannel(session, voiceState.GuildID)
	}
}

// generate a map of voice channel id:population
func getVCMembers(guildChannels []*discordgo.Channel, voiceStates []*discordgo.VoiceState) map[string]int {
	memberCount := make(map[string]int)

	// populate the map with all current voice channels
	for _, channel := range guildChannels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			memberCount[channel.ID] = 0
		}
	}

	// for each current voice state (member in voice channel)
	// add to their channel's membership number
	for _, state := range voiceStates {
		memberCount[state.ChannelID]++
	}

	return memberCount
}

// make a new voice channel with a random name
func makeNewVoiceChannel(session *discordgo.Session, guildID string) {
	// load channel names json file
	var channelNames []string
	file, _ := ioutil.ReadFile("./channel_names.json")
	_ = json.Unmarshal(file, &channelNames)

	// chose a random channel name
	channelName := channelNames[rand.Intn(len(channelNames))]

	// create a new voice channel under the "VOICE CHANNELS" parent
	_, _ = session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildVoice,
		Bitrate:  128000,
		ParentID: "379276406326165518", // TODO: dynamically get this ID
	})

	log.Printf("newvoicechannel: added new channel '%s'\n", channelName)
}

// delete all voice channels and make new ones, to change names
// TODO: make this avoid populated channels
func shuffleVCs(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// delete all voice channels
	guild, _ := session.State.Guild(msg.GuildID)
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			_, err := session.ChannelDelete(channel.ID) //! deleting populated channels is causing error "Uknown Channel"
			if err != nil {
				log.Println("\u001b[31mERROR:\u001b[0m", err.Error())
			}
		}
	}

	log.Println("shuffle: cleared all channels")

	// remake the new channels, drawing new names
	for i := 0; i < BUFFER_CHANNELS; i++ {
		makeNewVoiceChannel(session, guild.ID)
	}
}
