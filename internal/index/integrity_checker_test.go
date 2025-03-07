package index_test

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LukasJenicek/ggit/internal/index"
)

func TestCheckIndexIntegrity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		content     []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty content",
			content:     []byte{},
			expectError: false,
		},
		{
			name:        "Content too short",
			content:     []byte("short"),
			expectError: true,
			errorMsg:    "index header not found",
		},
		{
			name:        "Invalid header",
			content:     append([]byte("WRNG"), make([]byte, 8)...),
			expectError: true,
			errorMsg:    "invalid header",
		},
		{
			name: "Version mismatch",
			content: func() []byte {
				content := append([]byte("DIRC"), make([]byte, 8)...)
				// Version 1 instead of 2
				binary.BigEndian.PutUint32(content[4:8], 1)

				return content
			}(),
			expectError: true,
			errorMsg:    "expected version '2'",
		},

		{
			name: "Checksum mismatch",
			content: func() []byte {
				content := append([]byte("DIRC"), make([]byte, 8)...)
				binary.BigEndian.PutUint32(content[4:8], 2)
				content = append(content, []byte("some content")...)
				checksum, _ := hex.DecodeString("0000000000000000000000000000000000000000")

				return append(content, checksum...)
			}(),
			expectError: true,
			errorMsg:    "checksum does not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := index.CheckIndexIntegrity(tt.content)

			if tt.expectError {
				require.Error(t, err)

				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
