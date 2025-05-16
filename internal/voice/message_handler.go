package voice

import (
	"sync"
	"time"

	"go-bot/internal/logging"
	"go-bot/internal/tts"

	"github.com/bwmarrin/discordgo"
)

var (
	messageQueue = make(map[string]chan *discordgo.MessageCreate)
	queueMutex   sync.Mutex
)

// ProcessMessage is called when a new message is received
func ProcessMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	// Log incoming message details
	guild, err := s.State.Guild(m.GuildID)
	guildName := "Unknown"
	if err == nil {
		guildName = guild.Name
	}

	channel, err := s.State.Channel(m.ChannelID)
	channelName := "Unknown"
	if err == nil {
		channelName = channel.Name
	}

	logging.Info("MSG: Guild: '%s' (%s), Channel: '#%s' (%s), User: '%s#%s' (%s): %s",
		guildName, m.GuildID,
		channelName, m.ChannelID,
		m.Author.Username, m.Author.Discriminator, m.Author.ID,
		m.Content)

	// Log to dedicated message log file
	logging.LogMessage(
		guildName, m.GuildID,
		channelName, m.ChannelID,
		m.Author.Username, m.Author.Discriminator, m.Author.ID,
		m.Content)

	// Check if queue exists for this guild
	queueMutex.Lock()
	if _, exists := messageQueue[m.GuildID]; !exists {
		messageQueue[m.GuildID] = make(chan *discordgo.MessageCreate, 100)
		go processQueue(s, m.GuildID)
	}
	queueMutex.Unlock()

	// Add message to queue
	messageQueue[m.GuildID] <- m
}

// Process messages from the queue
func processQueue(s *discordgo.Session, guildID string) {
	for m := range messageQueue[guildID] {
		logging.Debug("Processing message: %s", m.Content)
		handleMessage(s, m)
	}
}

// Handle an individual message
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if bot is in a voice channel
	vc, ok := GetVoiceConnection(m.GuildID)
	if !ok || vc == nil || vc.ChannelID == "" {
		return
	}

	// Skip messages that shouldn't be processed
	if shouldSkipMessage(m) {
		return
	}

	// Process message content if not empty
	if m.Content != "" {
		// Apply both standard and user-specific text transformations
		text := ProcessUserText(m.Content, m.Author.ID)
		logging.Info("TTS text: %s", text)

		// Get or create the skip signal channel for this guild
		skipChan := GetSkipChannel(m.GuildID)
		skipDetected := false

		// Start a goroutine to listen for skip signals during TTS playback
		go func() {
			select {
			case <-skipChan:
				skipDetected = true
				logging.Info("Skip signal received for guild %s", m.GuildID)
			case <-time.After(5 * time.Minute): // Safety timeout
				// Do nothing, just a safety measure
			}
		}()

		// Generate and play audio
		if err := tts.SpeakText(text, vc); err != nil {
			logging.Error("TTS failed: %v", err)
			return
		}

		if skipDetected {
			logging.Info("TTS playback was skipped")
		} else {
			logging.Debug("TTS playback completed normally")
		}

		// Delete the message that was read
		logging.Debug("Attempting to delete message...")
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			logging.Warn("Failed to delete message: %v", err)
		} else {
			logging.Debug("Message deleted successfully.")
		}
	}
}

// Determine if a message should be skipped
func shouldSkipMessage(m *discordgo.MessageCreate) bool {
	return ContainsSticker(len(m.StickerItems)) ||
		ContainsMention(m.Content) ||
		IsLink(m.Content)
}
