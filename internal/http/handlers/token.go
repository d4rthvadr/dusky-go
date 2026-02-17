package handlers

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashAndEncodeToken takes a plain token string and returns its SHA-256 hash as a hexadecimal string
func hashAndEncodeToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])

}
