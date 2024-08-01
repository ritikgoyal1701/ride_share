package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHashedPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hashedPassword
}
