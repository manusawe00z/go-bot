package voice

import (
	"fmt"
	"sync"

	"go-bot/internal/logging"
	"go-bot/internal/tts"
)

var (
	// Map to store signal channels for skipping TTS messages
	skipSignals = make(map[string]chan struct{})
	skipMutex   sync.RWMutex
)

// GetSkipChannel returns the skip channel for a guild, creating it if needed
func GetSkipChannel(guildID string) chan struct{} {
	skipMutex.Lock()
	defer skipMutex.Unlock()

	if _, exists := skipSignals[guildID]; !exists {
		skipSignals[guildID] = make(chan struct{}, 1)
	}

	return skipSignals[guildID]
}

// RemoveSkipChannel removes a skip channel for a guild
func RemoveSkipChannel(guildID string) {
	skipMutex.Lock()
	defer skipMutex.Unlock()

	if ch, exists := skipSignals[guildID]; exists {
		close(ch)
		delete(skipSignals, guildID)
	}
}

// SkipCurrentMessage signals to skip the current TTS message
func SkipCurrentMessage(guildID string) error {
	skipMutex.RLock()
	skipChan, exists := skipSignals[guildID]
	skipMutex.RUnlock()

	if !exists {
		return fmt.Errorf("no active TTS session to skip")
	}

	// Try to send skip signal without blocking
	select {
	case skipChan <- struct{}{}:
		logging.Info("Skipping current TTS message for guild %s", guildID)

		// Also stop any active audio playback
		if session, exists := tts.GetAudioSession(guildID); exists {
			// Signal to stop the audio playback
			close(session.StopChan)
			logging.Info("Stopped active audio playback")
		}

		return nil
	default:
		// Channel already has a signal or is closed
		return nil
	}
}
