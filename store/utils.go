package store

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/everFinance/goar/utils"
)

// intToBase64 returns an 64-byte big endian representation of v.
func intToBase64(v uint64) string {
	b := make([]byte, 64)
	binary.BigEndian.PutUint64(b, v)
	return utils.Base64Encode(b)
}

func base64ToInt(base64Str string) uint64 {
	b, err := utils.Base64Decode(base64Str)
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint64(b)
}

func generateOffsetKey(root, size string) string {
	hash := sha256.Sum256([]byte(root + size))
	return utils.Base64Encode(hash[:])
}
