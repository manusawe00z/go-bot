package voice

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	voiceConnections = make(map[string]*discordgo.VoiceConnection)
	voiceMutex       sync.RWMutex
)

// SetVoiceConnection stores a voice connection for a guild
func SetVoiceConnection(guildID string, vc *discordgo.VoiceConnection) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()
	voiceConnections[guildID] = vc
}

// GetVoiceConnection retrieves a voice connection for a guild
func GetVoiceConnection(guildID string) (*discordgo.VoiceConnection, bool) {
	voiceMutex.RLock()
	defer voiceMutex.RUnlock()
	vc, ok := voiceConnections[guildID]
	return vc, ok
}

// RemoveVoiceConnection removes a voice connection for a guild
func RemoveVoiceConnection(guildID string) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()
	delete(voiceConnections, guildID)
}

// JoinUserVoice joins the voice channel of a user
func JoinUserVoice(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceConnection, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			// Set mute and deaf to true to disable bot input/output
			return s.ChannelVoiceJoin(guildID, vs.ChannelID, true, true)
		}
	}

	return nil, fmt.Errorf("user is not in a voice channel")
}
