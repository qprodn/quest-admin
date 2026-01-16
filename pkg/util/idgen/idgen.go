package idgen

import (
	"crypto/rand"
)

// GenerateID generates a unique ID
func GenerateID() string {
	return rand.Text()
}
