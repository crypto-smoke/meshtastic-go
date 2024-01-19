package mqtt

import (
	"math/rand"
	"strings"
	"time"
)

// generates a random string for use as a client ID
func randomString(n int, alphabet []rune) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[seededRand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}
