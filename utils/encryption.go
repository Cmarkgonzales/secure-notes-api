package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

var encryptionKey []byte

func LoadEncryptionKey() {
	raw := os.Getenv("ENCRYPTION_KEY")
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		panic("Invalid ENCRYPTION_KEY: must be base64 encoded")
	}
	if len(key) != 32 {
		panic("Encryption key must be 32 bytes for AES-256")
	}
	encryptionKey = key
	fmt.Println("Encryption key loaded successfully")
}

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(encoded string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(cipherText) < gcm.NonceSize() {
		return "", errors.New("cipherText too short")
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherData := cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
