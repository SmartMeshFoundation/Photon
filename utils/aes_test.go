package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key, _ := hex.DecodeString(passkey)
	// encript
	encodeBytes, err := Encrypt([]byte("some plaintext"), key)
	if err != nil {
		t.Error("Encrypt: ", err)
		return
	}
	// encode println
	fmt.Printf("Encrypt code: %x\n", string(encodeBytes))

	// decrypt
	decodeBytes, err := Decrypt(encodeBytes, key)
	if err != nil {
		t.Error("Decrypt: ", err)
		return
	}

	// decode println
	fmt.Println("Decrypt code: ", string(decodeBytes))
}

func TestPasswordDecrypt(t *testing.T) {
	encstr, err := PasswordEncrypt("1235566")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("encpass=%s", encstr)
	pass, err := PasswordDecrypt(encstr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("plain=%s", pass)
}
