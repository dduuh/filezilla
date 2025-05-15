package hash

import (
	"crypto/sha1"
	"fmt"
)

const salt = "gopher"

func HashPassword(password string) string {
	hash := sha1.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
