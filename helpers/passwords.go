package helpers

import (
	"crypto/rand"
	"fmt"
	"unsafe"
)

var alphanumeric = []byte("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ0123456789")

func GeneratePassword(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error occurred when generating password: %w", err)
	}
	for i := 0; i < length; i++ {
		b[i] = alphanumeric[b[i]%byte(len(alphanumeric))]
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}
