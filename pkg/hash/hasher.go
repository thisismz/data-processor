package hash

import (
	"fmt"

	"github.com/zeebo/blake3"
)

func Blake3Hash(payload []byte) string {
	h := blake3.New()
	h.Write(payload)
	return fmt.Sprintf("%x\n", h.Sum(nil))
}
