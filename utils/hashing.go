package utils

import (
	"crypto/sha256"
	"fmt"
)

func Hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashBytes := h.Sum(nil)

	return fmt.Sprintf("%x", hashBytes)
}
