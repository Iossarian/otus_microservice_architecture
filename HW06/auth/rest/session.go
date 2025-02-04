package rest

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

// sessions хранит соответствие sessionID -> userID
var sessions = struct {
	sync.RWMutex
	m map[string]int // sessionID -> userID
}{m: make(map[string]int)}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
