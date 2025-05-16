package bot

import (
	"context"
	"fmt"
	"time"

	"go-bot/internal/commands"
	"go-bot/internal/config"
	"go-bot/internal/http"
	"go-bot/internal/logging"
	"go-bot/internal/voice"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the Discord bot
type Bot struct {
	Session    *discordgo.Session
	Config     *config.Config
	HttpServer *http.Server
}

// NewBot creates a new bot instance
func NewBot() (*Bot, error) {
	// Load configuration
	cfg := config.LoadConfig()
	// Initialize logging
	if err := logging.Setup(cfg.LogLevel); err != nil {
		return nil, fmt.Errorf("failed to set up logging: %w", err)
	}

	// Initialize message logging
	if err := logging.SetupMessageLogger(); err != nil {
		return nil, fmt.Errorf("failed to set up message logging: %w", err)
	}

	// Validate Discord token
	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable not set")
	}
	// Create Discord session
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Enable the message content intent - required to read message contents
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentsMessageContent

	// Create HTTP server for health checks
	httpServer := http.NewServer(8080)

	// Create bot
	bot := &Bot{
		Session:    session,
		Config:     cfg,
		HttpServer: httpServer,
	}
	// Add message handler
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Skip messages from self
		if m.Author.ID == s.State.User.ID {
			return
		}
		// Log incoming message
		logging.Info("MSG: Guild: (%s), Channel: (%s), User: '%s#%s' (%s): %s",
			m.GuildID, m.ChannelID,
			m.Author.Username, m.Author.Discriminator, m.Author.ID,
			m.Content)

		voice.ProcessMessage(s, m)
	})

	// Add interaction handler for slash commands
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.handleInteraction(s, i)
	})

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	// Start HTTP server for health checks
	b.HttpServer.Start()

	// Open connection to Discord
	logging.Info("Starting bot...")
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	// Clear and register commands
	if err := commands.ClearCommands(b.Session); err != nil {
		logging.Warn("Failed to clear commands: %v", err)
	}

	if err := commands.RegisterCommands(b.Session); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	logging.Info("Bot is now running")
	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	logging.Info("Shutting down bot...")

	// Close Discord connection
	b.Session.Close()

	// Gracefully shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.HttpServer.Stop(ctx); err != nil {
		logging.Error("Error shutting down HTTP server: %v", err)
	}

	logging.CloseMessageLogger()
	logging.Close()
}

// handleInteraction processes slash commands
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Skip non-command interactions
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	logging.Info("Received command: %s", commandName)

	// Find and execute command
	for _, cmd := range commands.GetCommands() {
		if cmd.Name == commandName {
			if err := cmd.Handler(s, i); err != nil {
				logging.Error("Error handling command %s: %v", commandName, err)

				// Respond with error message
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Error: %v", err),
					},
				})
			}
			return
		}
	}

	// Command not found
	logging.Warn("Unknown command: %s", commandName)
}
