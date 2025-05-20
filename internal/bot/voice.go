package bot

import (
	"fmt"
	"os/exec"
	"regexp" // เพิ่ม import สำหรับ regex
	"strings"
	"sync"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	voiceConnections = make(map[string]*discordgo.VoiceConnection)
	messageQueue     = make(map[string]chan *discordgo.MessageCreate) // คิวข้อความสำหรับแต่ละเซิร์ฟเวอร์
	queueMutex       sync.Mutex                                       // ป้องกันการเข้าถึง map พร้อมกัน
)

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// ตรวจสอบว่ามีคิวสำหรับเซิร์ฟเวอร์นี้หรือยัง
	queueMutex.Lock()
	if _, exists := messageQueue[m.GuildID]; !exists {
		messageQueue[m.GuildID] = make(chan *discordgo.MessageCreate, 100) // สร้างคิวใหม่
		go b.processQueue(s, m.GuildID)                                    // เริ่มประมวลผลคิว
	}
	queueMutex.Unlock()

	// เพิ่มข้อความลงในคิว
	messageQueue[m.GuildID] <- m
}

func (b *Bot) processQueue(s *discordgo.Session, guildID string) {
	for m := range messageQueue[guildID] {
		fmt.Printf("Processing message: %s\n", m.Content)
		b.handleMessage(s, m)
	}
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ถ้าบอทยังไม่ได้ join ห้องเสียง
	vc, ok := voiceConnections[m.GuildID]
	if !ok || vc == nil || vc.ChannelID == "" {
		return
	}

	// ตรวจสอบข้อความที่ไม่ต้องประมวลผล เช่น ลิงก์, การแท็กบอท, หรือสติกเกอร์
	if len(m.StickerItems) > 0 || strings.Contains(m.Content, "@") || regexp.MustCompile(`https?://[^\s]+`).MatchString(m.Content) {
		return
	}

	// อ่านข้อความในแชนแนลและพูด
	if m.Content != "" {
		text := m.Content
		responseMuklock, isMuklock := muklock(text)
		fmt.Println("TTS text:", responseMuklock)
		// สร้างไฟล์เสียงจาก gTTS
		cmd := exec.Command("python3", "internal/tts/tts.py", text)
		err := cmd.Run()
		if err != nil {
			fmt.Println("TTS failed:", err)
			return
		}
		if isMuklock {
			// สร้างไฟล์เสียงจาก gTTS
			cmd := exec.Command("python3", "internal/tts/response-muklock.py", responseMuklock)
			err := cmd.Run()
			if err != nil {
				fmt.Println("TTS failed:", err)
				return
			}
		}
		// เล่นเสียงและรอให้การเล่นจบ
		fmt.Println("Starting to play audio...")

		if isMuklock {
			// Then play the original TTS
			go func() {
				done := make(chan bool)
				dgvoice.PlayAudioFile(vc, "tts.mp3", done)
				dgvoice.PlayAudioFile(vc, "response-muklock.mp3", done)
				dgvoice.PlayAudioFile(vc, "tlktbmuk.mp3", done)
				<-done
			}()

			// Finally play the muklock sound effect
		} else {
			// Play normal TTS
			go func() {
				done := make(chan bool)
				dgvoice.PlayAudioFile(vc, "tts.mp3", done)
				<-done
			}()
		}

		fmt.Println("Audio playback finished.")

		// ลบข้อความที่อ่านไปแล้ว
		fmt.Println("Attempting to delete message...")
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Printf("Failed to delete message: %v\n", err)
		} else {
			fmt.Println("Message deleted successfully.")
		}
	}
}

func (b *Bot) joinUserVoice(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceConnection, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			// ตั้งค่า mute และ deaf เป็น true เพื่อปิดเสียงเข้าและออกของบอท
			return s.ChannelVoiceJoin(guildID, vs.ChannelID, false, false)
		}
	}

	return nil, fmt.Errorf("user is not in a voice channel")
}

func muklock(text string) (string, bool) {
	switch text {
	case "สีเหลือง":
		return "เย็นโล่", true
	case "มีด":
		return "อีโต้", true
	case "ชุดชั้นใน":
		return "วาโก้", true
	}
	return text, false
}
