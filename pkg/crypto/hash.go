package crypto

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"strings"
)

func GenerateHash(data ...string) string {
	combinedString := strings.Join(data, "-")
	hash := sha256.New()

	hash.Write([]byte(combinedString))

	hashBytes := hash.Sum(nil)
	hashHex := fmt.Sprintf("%x", hashBytes)

	return hashHex
}

func GenerateLHash(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}
