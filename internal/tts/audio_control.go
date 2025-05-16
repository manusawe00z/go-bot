package tts

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	// Map of active audio playback sessions, indexed by guildID
	audioSessions = make(map[string]*AudioSession)
	sessionMutex  sync.RWMutex
)

// AudioSession represents an active audio playback
type AudioSession struct {
	// The channel that can be closed to stop playback
	StopChan chan bool
	// The voice connection being used
	VoiceConnection *discordgo.VoiceConnection
}

// GetAudioSession returns the active audio session for a guild
func GetAudioSession(guildID string) (*AudioSession, bool) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	session, exists := audioSessions[guildID]
	return session, exists
}

// CreateAudioSession creates a new audio session for a guild
func CreateAudioSession(guildID string, vc *discordgo.VoiceConnection) *AudioSession {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	// Close existing session if it exists
	if session, exists := audioSessions[guildID]; exists {
		close(session.StopChan)
	}

	// Create a new session
	session := &AudioSession{
		StopChan:        make(chan bool),
		VoiceConnection: vc,
	}

	audioSessions[guildID] = session
	return session
}

// RemoveAudioSession removes an audio session for a guild
func RemoveAudioSession(guildID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if session, exists := audioSessions[guildID]; exists {
		close(session.StopChan)
		delete(audioSessions, guildID)
	}
}
