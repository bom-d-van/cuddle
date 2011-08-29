package cuddle

import (
	"rand"
	"regexp"
	"time"
)

// ValidName matches a valid name string.
var ValidName = regexp.MustCompile(`^[a-z]+$`)

// RandName returns a string of l random characters.
func RandName(l int) string {
	n := make([]byte, l)
	for i := range n {
		n[i] = 'a' + byte(rand.Intn(26))
	}
	return string(n)
}

func init() {
	// Seed number generator with the current time.
	rand.Seed(time.Nanoseconds())
}
