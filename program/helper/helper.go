package helper

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func GravatarUrl(email string) string {
	size := 40
	h := sha1.New()
	h.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", sha1Hash, size)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}