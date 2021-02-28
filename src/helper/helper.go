package helper

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func Gravatar_url(email string) string {
	size := 40
	h := sha1.New()
	h.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", sha1_hash, size)
}
