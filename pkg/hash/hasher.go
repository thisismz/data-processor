package hash

import "github.com/zeebo/blake3"

func Blake3Hash(payload []byte) string {
	h := blake3.New()
	h.Write(payload)
	return string(h.Sum(nil))
}
