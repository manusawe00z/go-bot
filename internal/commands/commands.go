package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Command represents a bot command
type Command struct {
	Name        string
	Description string
	Handler     CommandHandler
}

// CommandHandler is a function that handles a command
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

// Commands returns all available commands
func GetCommands() []*Command {
	return []*Command{
		{
			Name:        "join",
			Description: "Join the voice channel",
			Handler:     HandleJoin,
		},
		{
			Name:        "leave",
			Description: "Leave the voice channel",
			Handler:     HandleLeave,
		},
		{
			Name:        "skip",
			Description: "Skip the current TTS message",
			Handler:     HandleSkip,
		},
		{
			Name:        "hate",
			Description: "Add a user to the special TTS transformation list",
			Handler:     HandleHate,
		},
		{
			Name:        "unhate",
			Description: "Remove a user from the special TTS transformation list",
			Handler:     HandleUnhate,
		},
		{
			Name:        "hatelist",
			Description: "Show all users in the special TTS transformation list",
			Handler:     HandleHateList,
		},
	}
}

// RegisterCommands registers all slash commands for a bot
func RegisterCommands(s *discordgo.Session) error {
	commands := GetCommands()
	discordCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, cmd := range commands {
		// Create the base command
		discordCommand := &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		}

		// Add command-specific options
		switch cmd.Name {
		case "hate", "unhate":
			// Add a required user option for the hate/unhate commands
			discordCommand.Options = []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to target",
					Required:    true,
				},
			}
		}

		discordCommands[i] = discordCommand
	}

	// Register commands for each guild
	for _, guild := range s.State.Guilds {
		fmt.Printf("Registering commands for guild: %s (%s)\n", guild.Name, guild.ID)
		for _, cmd := range discordCommands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, cmd)
			if err != nil {
				return fmt.Errorf("cannot create '%s' command for guild %s: %v", cmd.Name, guild.Name, err)
			}
		}
	}

	return nil
}

// ClearCommands removes all registered commands
func ClearCommands(s *discordgo.Session) error {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("cannot fetch commands: %v", err)
	}

	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			return fmt.Errorf("cannot delete '%s' command: %v", cmd.Name, err)
		}
	}

	fmt.Println("All commands cleared.")
	return nil
}
