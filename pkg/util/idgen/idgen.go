package idgen

import (
	"crypto/rand"
)

func init() {

}

// GenerateID generates a unique ID
func GenerateID() string {
	return rand.Text()
}
