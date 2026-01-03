package helpers

import "crypto/sha256"

func calculateSHA256(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
