package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	separator     = "--"
	authTagLength = 16
)

// Decrypt decrypts the input using AES-GCM
func Decrypt(encData string, key []byte) ([]byte, error) {
	encParts := strings.Split(encData, separator)
	if len(encParts) != 3 {
		return nil, errors.New("invalid encrypted format")
	}

	encryptedData, err := base64.StdEncoding.Strict().DecodeString(encParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid base64 for encrypted encData: %w", err)
	}

	iv, err := base64.StdEncoding.Strict().DecodeString(encParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid base64 for iv: %w", err)
	}

	authTag, err := base64.StdEncoding.Strict().DecodeString(encParts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid base64 for auth_tag: %w", err)
	}

	if len(authTag) != authTagLength {
		return nil, errors.New("wrong auth tag length")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to init AES: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to init GCM: %w", err)
	}

	fullCiphertext := append(encryptedData, authTag...)

	return gcm.Open(nil, iv, fullCiphertext, nil)
}

// Encrypt encrypts the input using AES-GCM
func Encrypt(data []byte, key []byte) (string, error) {
	if len(key) != 16 && len(key) != 32 {
		return "", fmt.Errorf("unsupported key length: %d", len(key))
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to init AES: %w", err)
	}

	gcm, err := cipher.NewGCMWithTagSize(aesCipher, authTagLength)
	if err != nil {
		return "", fmt.Errorf("failed to init GCM: %w", err)
	}

	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	input := data
	ciphertext := gcm.Seal(nil, iv, input, nil)
	authTag := ciphertext[len(ciphertext)-authTagLength:]
	encData := ciphertext[:len(ciphertext)-authTagLength]

	parts := [][]byte{encData, iv, authTag}
	var encoded []string
	for _, part := range parts {
		encoded = append(encoded, base64.StdEncoding.Strict().EncodeToString(part))
	}

	return strings.Join(encoded, separator), nil
}
