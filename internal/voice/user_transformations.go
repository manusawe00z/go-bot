package voice

import (
	"math/rand"
	"sync"
	"time"

	"go-bot/internal/logging"
)

func init() {
	// In Go 1.20+ we don't need to explicitly seed the random number generator
	// This is kept for compatibility with older Go versions
	rand.Seed(time.Now().UnixNano())
}

var (
	// Map of user IDs to whether they are "hated" (targeted for special transformations)
	hatedUsers = make(map[string]bool)
	// Mutex for thread-safe access to the hatedUsers map
	hatedUsersMutex sync.RWMutex
)

// AddHatedUser adds a user to the "hated" list for special text transformations
func AddHatedUser(userID string) {
	hatedUsersMutex.Lock()
	defer hatedUsersMutex.Unlock()

	hatedUsers[userID] = true
	logging.Info("Added user %s to hated users list", userID)
}

// RemoveHatedUser removes a user from the "hated" list
func RemoveHatedUser(userID string) {
	hatedUsersMutex.Lock()
	defer hatedUsersMutex.Unlock()

	delete(hatedUsers, userID)
	logging.Info("Removed user %s from hated users list", userID)
}

// IsHatedUser checks if a user is in the "hated" list
func IsHatedUser(userID string) bool {
	hatedUsersMutex.RLock()
	defer hatedUsersMutex.RUnlock()

	return hatedUsers[userID]
}

// GetAllHatedUsers returns a list of all hated user IDs
func GetAllHatedUsers() []string {
	hatedUsersMutex.RLock()
	defer hatedUsersMutex.RUnlock()

	users := make([]string, 0, len(hatedUsers))
	for userID := range hatedUsers {
		users = append(users, userID)
	}

	return users
}

// ApplyUserSpecificTransformations applies transformations based on the user who sent the message
func ApplyUserSpecificTransformations(text string, userID string) string {
	// Check if this user is in the hated list
	if !IsHatedUser(userID) {
		return text
	}

	// Apply special transformations for hated users
	return ApplyHatedUserTransformations(text)
}

// ApplyHatedUserTransformations applies special transformations for hated users
func ApplyHatedUserTransformations(text string) string {
	// If the text is empty or we want to completely replace it
	if text == "" || rand.Intn(100) < 90 { // 10% chance to replace message completely
		// List of sarcastic responses
		responses := []string{
			"กลัวความสูงเหรอ เห็นชอบทำตัวต่ำๆ",
			"การศึกษาไม่ได้ทำให้คนฉลาดทางอารมณ์",
			"กินอาหารดีๆ บ้างนะ จะได้มีสารอาหารไปเลี้ยงสมอง",
			"เก็บปากไว้แตกหน้าหนาวเถอะ",
			"เก็บแรงด่า ไว้แผ่เมตตาให้ดีกว่า",
			"เกลียดนักพวกที่ชอบด่าคนอื่นไม่ดูตัวเอง",
			"เกิดปีไก่เหรอ จิกขนาดนี้",
			"คนบางคนก็เหมือนฝุ่น PM 2.5 ดูไร้ค่า และยังเป็นพิษกับคนอื่น",
			"คนหรือกระดาษทราย ดูตัวหยาบๆ",
			"เคยเป็นช่างไฟฟ้าเหรอ เห็นชอบสร้างกระแสเก่งจัง",
			"จะดูถูกกันไม่ว่า แต่ช่วยดูหน้าตัวเองด้วย",
			"ช่วงนี้เที่ยวทิพย์ได้ แต่อย่าเที่ยวยุ่งเรื่องคนอื่นนะ",
			"ชอบปลายปี ชอบหน้าหนาว แต่ไม่ชอบหน้าเธออะ",
			"ดูๆ ก็เหมือนจะครบ 32 ขาดแต่สมองอย่างเดียว",
			"ถ้าชอบกัดก็ไปเกิดเป็นสัตว์นะ",
			"ถึงไม่ได้ขายเพชรพลอย แต่ก็ดูออกว่าใครปลอม",
			"ที่บ้านไม่เคยสอนเหรอว่าอย่าไปอยากได้ของของคนอื่น",
			"ที่เห็นเงียบๆ นี่รอเหยียบนะ ไม่ได้ยอม",
			"เธอก็เหมือนแมงกะพรุน เห็นใสๆ แต่มีพิษ",
			"คนอย่างเธอมันก็เหมือนแบงก์กาโม่อ่ะ ถูกและปลอม",
			"นี่ภาคกลางจ้า ไม่ต้องโชว์เหนือ",
			"บางคำที่คุณพูดอาจไม่ใช่เรื่องตลกสำหรับคนอื่นเสมอไป",
			"ปกตินี่ถ้าหิวน้ำ ต้องกิน หรือต้องมีคนกรวดให้",
			"พร้อมบวกตลอดนะโทรศัพท์มีเครื่องคิดเลข",
			"พูดไปเถอะเรื่องของคนอื่น ตายไปก็พูดไม่ได้อยู่ดี",
			"มารยาททางสังคมไม่มี สงสัยบุพการีไม่ได้สอน",
			"ไม่รู้อะไรอย่าพูด เก็บน้ำลายไว้ย่อยอาหารยังจะมีประโยชน์กว่า",
			"เลิกกันแล้วอย่าเรียกแฟนเก่า อะไรที่ไม่เอา เขาเรียก 'ขยะ'",
			"ช่วงนี้ร้อนเงินเหรอคะ เห็นขยันลดราคาตัวเองจัง",
			"สมองเธอคงเรียบหรู เป็นคุณหนูแต่ไม่มีรอยหยัก",
			"สมองไม่มีไม่เป็นไร แต่มารยาทต้องมีนะคะ",
			"สร้างภาพเก่งขนาดนี้ มีอาชีพเป็นช่างภาพเหรอ",
			"หน้าก็ลอยแล้ว กระทงไม่ต้องลอยก็ได้มั้ง",
			"หน้าตาก็ดูฉลาด แต่พูดจาเหมือนขาดสมอง",
			"หน้าเธอเหมือนสี่เหลี่ยมจัตุรัสเลย ด้านคูณด้าน",
			"หน้านะไม่ใช่แบงก์พัน ที่จะได้เทาแล้วดูแพง",
			"หลอกเก่งจริงๆ ลืมไปหรือเปล่าว่าตัวเองยังไม่ตายนะ",
			"อยากเป็นส่วนหนึ่ง แต่ระวังจะเป็นส่วนเกิน",
			"อยากสูงขึ้นก็กินนมนะ ไม่ใช่เหยียบหัวคนอื่น",
			"อยู่ไหนเหรอความจริงใจ เราหาจากเธอเท่าไรก็หาไม่เจอ",
			"ไม่มีสมอง",
		}

		return responses[rand.Intn(len(responses))]
	}

	return text
}
