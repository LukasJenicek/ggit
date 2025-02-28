package index

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/LukasJenicek/ggit/internal/hasher"
)

func CheckIndexIntegrity(content []byte) error {
	if len(content) == 0 {
		return nil
	}

	if len(content) < 12 {
		return fmt.Errorf("index header not found")
	}

	header := content[:12]
	if string(header[:4]) != "DIRC" {
		return fmt.Errorf("invalid header")
	}

	if binary.BigEndian.Uint32(header[4:8]) != 2 {
		return fmt.Errorf("expected version '2'")
	}

	hashContent, err := hasher.SHA1HashContent(content[:len(content)-20])
	if err != nil {
		return fmt.Errorf("hash content: %w", err)
	}
	checksum := content[len(content)-20:]

	if hex.EncodeToString(checksum) != hex.EncodeToString(hashContent) {
		return fmt.Errorf("checksum does not match")
	}

	return nil
}
