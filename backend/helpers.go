
package main

import (
	"crypto/rand"
	"encoding/hex"
)

func generateLotNumber() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
