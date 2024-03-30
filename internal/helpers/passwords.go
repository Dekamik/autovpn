package helpers

import (
	"crypto/rand"
	"unsafe"
)

var alphanumeric = []byte("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ0123456789")

// GeneratePassword generates a random password consisting of the characters
// [a-zA-Z0-9]
func GeneratePassword(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		b[i] = alphanumeric[b[i]%byte(len(alphanumeric))]
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}
