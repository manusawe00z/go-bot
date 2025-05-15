package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot() (*Bot, error) {
	botToken := os.Getenv("DISCORD_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable not set")
	}
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{Session: session}
	session.AddHandler(bot.messageCreate)
	return bot, nil
}

func (b *Bot) Start() error {
	// เปิดการเชื่อมต่อกับ Discord Gateway
	err := b.Session.Open()
	if err != nil {
		return err
	}

	// ลบคำสั่งเก่าก่อนลงทะเบียนใหม่
	err = b.ClearCommands()
	if err != nil {
		return err
	}

	// ลงทะเบียน Slash Commands หลังจากเปิดการเชื่อมต่อสำเร็จ
	err = b.RegisterCommands()
	if err != nil {
		return err
	}

	b.Session.AddHandler(b.InteractionCreate)
	return nil
}

func (b *Bot) Stop() {
	b.Session.Close()
}

func (b *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "join",
			Description: "ให้บอทเข้าร่วมห้องเสียง",
		},
		{
			Name:        "leave",
			Description: "ให้บอทออกจากห้องเสียง",
		},
	}

	// ดึง Guild ID จากเซิร์ฟเวอร์ที่บอทเข้าร่วม
	for _, guild := range b.Session.State.Guilds {
		fmt.Printf("Registering commands for guild: %s (%s)\n", guild.Name, guild.ID)
		for _, cmd := range commands {
			_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, guild.ID, cmd)
			if err != nil {
				return fmt.Errorf("cannot create '%s' command for guild %s: %v", cmd.Name, guild.Name, err)
			}
		}
	}

	return nil
}

func (b *Bot) ClearCommands() error {
	commands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("cannot fetch commands: %v", err)
	}

	for _, cmd := range commands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID)
		if err != nil {
			return fmt.Errorf("cannot delete '%s' command: %v", cmd.Name, err)
		}
	}

	fmt.Println("All commands cleared.")
	return nil
}

func (b *Bot) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "join":
		err := b.handleJoinCommand(s, i)
		if err != nil {
			fmt.Println("Error handling join command:", err)
		}
	case "leave":
		err := b.handleLeaveCommand(s, i)
		if err != nil {
			fmt.Println("Error handling leave command:", err)
		}
	}
}

func (b *Bot) handleJoinCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	vc, err := b.joinUserVoice(s, i.GuildID, i.Member.User.ID)
	if err != nil {
		return err
	}
	voiceConnections[i.GuildID] = vc

	// ตอบกลับคำสั่ง
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Joined the voice channel!",
		},
	})
}

func (b *Bot) handleLeaveCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	vc, ok := voiceConnections[i.GuildID]
	if !ok || vc == nil || vc.ChannelID == "" {
		return fmt.Errorf("bot is not in a voice channel")
	}

	err := vc.Disconnect()
	if err != nil {
		return err
	}

	delete(voiceConnections, i.GuildID)

	// ตอบกลับคำสั่ง
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Left the voice channel!",
		},
	})
}
