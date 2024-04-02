package utils_test

import (
	"fmt"
	"launchpad/utils"
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := utils.SBuilder.GetString(16)
		t.Log(s)
	}
}

func TestRc4(t *testing.T) {
	rc4 := utils.EncryptionRc4("api_key", "this is q")
	fmt.Println(rc4)
	decryptionRc4 := utils.DecryptionRc4("api_key", rc4)
	fmt.Println(decryptionRc4)
}
