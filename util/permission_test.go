package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedPermission(t *testing.T) {
	testCases := []struct {
		name     string
		perm     string
		expected bool
	}{
		{
			name:     "read",
			expected: true,
		},
		{
			name:     "write",
			expected: true,
		},
		{
			name:     "delete",
			expected: true,
		},
		{
			name:     "unknown",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsSupportedPermission(tc.name)
			require.Equal(t, actual, tc.expected)
		})
	}
}
