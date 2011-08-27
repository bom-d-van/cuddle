package cuddle

import (
	"rand"
	"regexp"
	"time"
)

func init() {
	// Seed the random number generator with the current time.
	rand.Seed(time.Nanoseconds())
}

// nameChars are the characters that can be in a name.
const nameChars = "abcdefghijklmnopqrstuvwxyz"

// ValidName matches a string consisting of nameChars.
var ValidName = regexp.MustCompile(`^[` + nameChars + `]+$`)

// RandName returns a string of l random characters chosen from nameChars.
func RandName(l int) string {
	n := make([]byte, l)
	for i := range n {
		n[i] = nameChars[rand.Intn(len(nameChars))]
	}
	return string(n)
}
