package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// randomString 產生隨機字串
func randomString(length int) string {
	bytes := make([]byte, length/2+1)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// generateConnectionID 產生連線 ID
func generateConnectionID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// generatePresetID 產生預設值 ID
func generatePresetID() string {
	return "preset-" + randomString(8)
}

// generateTemplateID 產生模板 ID
func generateTemplateID() string {
	return "template-" + randomString(8)
}
