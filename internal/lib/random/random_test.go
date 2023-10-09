package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "Test for length 1",
			size: 1,
		},
		{
			name: "Test for length 10",
			size: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			assert.NotEqualf(t, str1, str2,
				"The two random generated strings is equal: first: \"%s\", \"%s\"", str1, str2)
		})
	}
}
