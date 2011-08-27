package cuddle

import (
	"rand"
	"regexp"
	"time"
)

func init() {
	rand.Seed(time.Nanoseconds())
}

const nameChars = "abcdefghijklmnopqrstuvwxyz"

var validName = regexp.MustCompile(`^[` + nameChars + `]+$`)

func randName(l int) string {
	n := make([]byte, l)
	for i := range n {
		n[i] = nameChars[rand.Intn(len(nameChars))]
	}
	return string(n)
}
