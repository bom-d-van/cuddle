package cuddle

import (
	"rand"
	"regexp"
	"time"
)

func init() {
	rand.Seed(time.Nanoseconds())
}

const (
	nameChars = "abcdefghijklmnopqrstuv"
	nameLen   = 4
)

var validName = regexp.MustCompile(`^[` + nameChars + `]+$`)

func randName() string {
	n := make([]byte, nameLen)
	for i := range n {
		n[i] = nameChars[rand.Intn(len(nameChars))]
	}
	return string(n)
}
