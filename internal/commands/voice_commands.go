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
