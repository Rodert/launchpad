package utils

import (
	"crypto/rc4"
	"encoding/base64"
)

func EncryptionRc4(k, query string) string {
	key := []byte(k)
	plaintext := []byte(query)
	// encryption
	ciphertext := make([]byte, len(plaintext))
	cipher1, _ := rc4.NewCipher(key)
	cipher1.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptionRc4(k, query string) string {
	param, err := base64.StdEncoding.DecodeString(query)
	if err != nil {
		return ""
	}
	key := []byte(k)
	ciphertext := param
	plaintextDec := make([]byte, len(ciphertext))
	cipher2, _ := rc4.NewCipher(key)
	cipher2.XORKeyStream(plaintextDec, ciphertext)
	return string(plaintextDec)
}
