package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateTOTP generates a TOTP code using the secret
func GenerateTOTP(secret string) (string, error) {
	decodedSecret, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	counter := time.Now().Unix() / 30 // TOTP changes every 30 seconds
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	h := hmac.New(sha256.New, decodedSecret)
	h.Write(buf)

	hash := h.Sum(nil)
	offset := hash[len(hash)-1] & 0x0f
	code := binary.BigEndian.Uint32(hash[offset : offset+4])
	return fmt.Sprintf("%06d", code%1000000), nil
}
