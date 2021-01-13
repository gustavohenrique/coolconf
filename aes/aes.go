package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/scrypt"
)

// 32 bits for AES 256
const IV_SIZE = 32

func Encrypt(secret string, plaintext []byte) (string, error) {
	key, salt, err := deriveKey([]byte(secret), nil)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := randomKey(gcm.NonceSize())
	// prepend nonce into encrypted data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	ciphertext = append(ciphertext, salt...)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(secret string, data []byte) (string, error) {
	salt, ciphertext := data[len(data)-IV_SIZE:], data[:len(data)-IV_SIZE]
	key, _, err := deriveKey([]byte(secret), salt)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("The ciphertext length is short than the nonce size")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}

func deriveKey(secret, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = randomKey(IV_SIZE)
	}
	key, err := scrypt.Key(secret, salt, 1048576, 8, 1, IV_SIZE)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

func randomKey(size int) []byte {
	key := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		log.Fatalln(err)
	}
	return key
}
