package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/LukasJenicek/ggit/internal/hasher"
)

func CheckIndexIntegrity(content []byte) error {
	if len(content) == 0 {
		return nil
	}

	if len(content) < 12 {
		return errors.New("index header not found")
	}

	header := content[:12]
	if string(header[:4]) != "DIRC" {
		return errors.New("invalid header")
	}

	if binary.BigEndian.Uint32(header[4:8]) != 2 {
		return errors.New("expected version '2'")
	}

	hashContent, err := hasher.SHA1HashContent(content[:len(content)-20])
	if err != nil {
		return fmt.Errorf("hash content: %w", err)
	}

	checksum := content[len(content)-20:]

	if !bytes.Equal(checksum, hashContent) {
		return errors.New("checksum does not match")
	}

	return nil
}
