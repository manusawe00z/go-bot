package tts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go-bot/internal/config"
	"go-bot/internal/logging"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// TTS contains configuration for text-to-speech
type TTS struct {
	Language     string
	Voice        string
	OutputDir    string
	PythonPath   string
	ScriptPath   string
	AudioQuality string
}

// NewTTS creates a new TTS instance with configuration
func NewTTS(cfg *config.Config) *TTS {
	// Create output directory if it doesn't exist
	outputDir := filepath.Join(".", "audio")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logging.Warn("Failed to create audio directory: %v", err)
	}

	// Find Python executable
	pythonPath := "python3"
	if _, err := exec.LookPath(pythonPath); err != nil {
		pythonPath = "python"
		if _, err := exec.LookPath(pythonPath); err != nil {
			logging.Warn("Python executable not found, using default 'python3'")
			pythonPath = "python3"
		}
	}

	return &TTS{
		Language:     cfg.TTSLanguage,
		Voice:        cfg.TTSVoice,
		OutputDir:    outputDir,
		PythonPath:   pythonPath,
		ScriptPath:   "internal/tts/tts.py",
		AudioQuality: cfg.AudioQuality,
	}
}

// SpeakText converts text to speech and plays it in the voice channel
func SpeakText(text string, vc *discordgo.VoiceConnection) error {
	// Get language from environment or use default
	language := os.Getenv("TTS_LANGUAGE")
	if language == "" {
		language = "en"
	}

	// Generate a unique filename based on text content
	outputFile := "tts.mp3"

	logging.Debug("Generating TTS audio for text: %s", text)

	// Create the TTS command
	cmd := exec.Command("python3", "internal/tts/tts.py", text, language, outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate TTS audio: %w (output: %s)", err, string(output))
	}

	// Play the audio file
	logging.Info("Starting to play audio...")
	done := make(chan bool)
	go func() {
		defer close(done)
		dgvoice.PlayAudioFile(vc, outputFile, done)
	}()
	<-done
	logging.Info("Audio playback finished.")

	return nil
}
