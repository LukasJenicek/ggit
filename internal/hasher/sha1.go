package hasher

import (
	"crypto/sha1"
	"fmt"
)

func SHA1HashContent(content []byte) ([]byte, error) {
	hasher := sha1.New()
	if _, err := hasher.Write(content); err != nil {
		return nil, fmt.Errorf("sha1 hash content: %w", err)
	}

	return hasher.Sum(nil), nil
}
