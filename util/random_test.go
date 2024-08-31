package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	randomInt := randomInt(0, 10)
	require.Less(t, randomInt, int64(10))
	require.GreaterOrEqual(t, randomInt, int64(0))

	randomString := randomString(10)
	require.Equal(t, len(randomString), 10)

	randomEmail := RandomEmail()
	require.Contains(t, randomEmail, "@email.com")

	randomPassword := RandomPassword()
	require.NotEmpty(t, randomPassword)

	randomSpaceName := RandomSpaceName()
	require.Equal(t, len(randomSpaceName), 6)

}
