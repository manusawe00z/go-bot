package commands

import (
	"fmt"

	"go-bot/internal/voice"

	"github.com/bwmarrin/discordgo"
)

// HandleJoin handles the join command
func HandleJoin(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	vc, err := voice.JoinUserVoice(s, i.GuildID, i.Member.User.ID)
	if err != nil {
		return err
	}
	voice.SetVoiceConnection(i.GuildID, vc)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Joined the voice channel!",
		},
	})
}

// HandleLeave handles the leave command
func HandleLeave(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	vc, ok := voice.GetVoiceConnection(i.GuildID)
	if !ok || vc == nil || vc.ChannelID == "" {
		return fmt.Errorf("bot is not in a voice channel")
	}

	err := vc.Disconnect()
	if err != nil {
		return err
	}

	voice.RemoveVoiceConnection(i.GuildID)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Left the voice channel!",
		},
	})
}

// HandleSkip handles the skip command
func HandleSkip(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Try to skip current message
	err := voice.SkipCurrentMessage(i.GuildID)

	var responseContent string
	if err != nil {
		responseContent = fmt.Sprintf("Error: %v", err)
	} else {
		responseContent = "Skipped current message!"
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	})
}

// HandleHate handles the hate command to target specific users for text transformations
func HandleHate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the mentioned user from the options
	options := i.ApplicationCommandData().Options

	if len(options) != 1 || options[0].Type != discordgo.ApplicationCommandOptionUser {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: Please mention a user with @username",
			},
		})
	}

	// Get the targeted user
	targetUser := options[0].UserValue(s)
	if targetUser == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: Invalid user",
			},
		})
	}

	// Add the user to the "hated" list
	voice.AddHatedUser(targetUser.ID)

	// Respond to the interaction
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Added %s to the special TTS transformation list 😈", targetUser.Username),
		},
	})
}

// HandleUnhate handles the unhate command to remove a user from the special transformation list
func HandleUnhate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the mentioned user from the options
	options := i.ApplicationCommandData().Options

	if len(options) != 1 || options[0].Type != discordgo.ApplicationCommandOptionUser {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: Please mention a user with @username",
			},
		})
	}

	// Get the targeted user
	targetUser := options[0].UserValue(s)
	if targetUser == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: Invalid user",
			},
		})
	}

	// Remove the user from the "hated" list
	voice.RemoveHatedUser(targetUser.ID)

	// Respond to the interaction
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Removed %s from the special TTS transformation list", targetUser.Username),
		},
	})
}

// HandleHateList handles the hatelist command to show all users in the special transformation list
func HandleHateList(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get all hated users
	hatedUserIDs := voice.GetAllHatedUsers()

	content := "Special TTS transformation list:\n"

	if len(hatedUserIDs) == 0 {
		content = "No users in the special TTS transformation list."
	} else {
		// Lookup user info for each ID
		for _, userID := range hatedUserIDs {
			// Try to get user information
			user, err := s.User(userID)
			if err != nil {
				content += fmt.Sprintf("- Unknown user (ID: %s)\n", userID)
			} else {
				content += fmt.Sprintf("- %s#%s (ID: %s)\n", user.Username, user.Discriminator, user.ID)
			}
		}
	}

	// Respond to the interaction
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
