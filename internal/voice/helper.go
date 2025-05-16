package voice

import (
	"regexp"
	"strings"
)

// IsLink checks if a string contains a URL
func IsLink(text string) bool {
	return regexp.MustCompile(`https?://[^\s]+`).MatchString(text)
}

// ContainsMention checks if a string contains a mention
func ContainsMention(text string) bool {
	return strings.Contains(text, "@")
}

// ContainsSticker checks if a message has stickers
func ContainsSticker(stickerCount int) bool {
	return stickerCount > 0
}

// ProcessText applies text transformations before TTS
func ProcessText(text string) string {
	// Custom pronunciations/replacements
	switch text {
	case "สีเหลือง":
		return "เย็นโล่"
	case "มีด":
		return "อีโต้"
	case "ชุดชั้นใน":
		return "วาโก้"
	}

	// Add more transformations here as needed

	return text
}

// ProcessUserText applies text transformations based on user and message content
func ProcessUserText(text string, userID string) string {
	// First apply standard text transformations
	transformedText := ProcessText(text)

	// Then apply user-specific transformations if applicable
	return ApplyUserSpecificTransformations(transformedText, userID)
}
