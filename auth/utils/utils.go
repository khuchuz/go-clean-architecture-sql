package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashThis(plaintext, salt string) string {
	pwd := sha1.New()
	pwd.Write([]byte(plaintext))
	pwd.Write([]byte(salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
