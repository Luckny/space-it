package util

import (
	"strings"

	"golang.org/x/exp/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// generates a random integer between min and max
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 to max-min
}

// generates a random string of length n
func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generates a random email
func RandomEmail() string {
	return randomString(3) + "@" + randomString(4) + ".com"
}

// generates a random password
func RandomPassword() string {
	return randomString(10) + string(randomInt(0, 10))
}
