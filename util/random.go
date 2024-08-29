package util

import (
	"math/rand"
	"strconv"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generates a random integer between min and max
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 to max-min
}

// generates a random string of length n
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(b)
}

// generates a random email
func RandomEmail() string {
	return randomString(4) + "@email.com"
}

// generates a random password
func RandomPassword() string {
	return randomString(10) + strconv.FormatInt(randomInt(0, 10), 10)
}

// generates a random space name
func RandomSpaceName() string {
	return randomString(6)
}
