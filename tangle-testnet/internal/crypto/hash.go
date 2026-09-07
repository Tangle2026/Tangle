package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Hash(parts ...any) string {
	h := sha256.New()
	for _, p := range parts {
		fmt.Fprintf(h, "%v|", p)
	}
	return hex.EncodeToString(h.Sum(nil))
}
