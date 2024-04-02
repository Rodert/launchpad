package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetEncrypt(query string) string {
	return SHA1(query)
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}
