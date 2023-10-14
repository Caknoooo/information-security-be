package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
)

// Encrypt data using DES with PKCS7 padding
func DESEncrypt(stringToEncrypt string, KEYS string) (encryptedString string, data map[string]interface{}, err error) {
	key, _ := hex.DecodeString(KEYS)
	plaintext := []byte(stringToEncrypt)

	if len(key) != 8 {
		return "", nil, errors.New("kunci DES harus memiliki panjang 8 byte")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	// PKCS7 padding
	plaintext = PKCS7Pad(plaintext, des.BlockSize)

	mode := cipher.NewCBCEncrypter(block, make([]byte, des.BlockSize))

	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	data = map[string]interface{}{
		"key":         KEYS,
		"plaintext":   string(plaintext),
		"block":       fmt.Sprintf("%d", block.BlockSize()),
		"mode_chiper": fmt.Sprintf("%v", mode),
		"mode":        "DES",
	}

	return fmt.Sprintf("%x", ciphertext), data, err
}

// Decrypt data using DES with PKCS7 padding
func DESDecrypt(encryptedString string, KEYS string) (decryptedString string, err error) {
	key, _ := hex.DecodeString(KEYS)
	enc, _ := hex.DecodeString(encryptedString)

	if len(key) != 8 {
		return "", errors.New("kunci DES harus memiliki panjang 8 byte")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, make([]byte, des.BlockSize))

	plaintext := make([]byte, len(enc))
	mode.CryptBlocks(plaintext, enc)

	// Remove PKCS7 padding
	plaintext, err = PKCS7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// PKCS7Pad adds PKCS7 padding to the data to make it a multiple of blockSize
func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7Unpad removes PKCS7 padding from the data
func PKCS7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid padding")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}
