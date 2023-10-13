package utils

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
)

// DESEncrypt melakukan enkripsi DES
func DESEncrypt(stringToEncrypt string, KEYS string) (encryptedString string, data map[string]interface{}, err error) {
	key, _ := hex.DecodeString(KEYS)
	plaintext := []byte(stringToEncrypt)

	// Pastikan panjang kunci DES adalah 8 byte
	if len(key) != 8 {
		return "", nil, errors.New("kunci DES harus memiliki panjang 8 byte")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	// Mode ECB untuk enkripsi DES
	mode := cipher.NewCBCEncrypter(block, make([]byte, des.BlockSize))

	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	data = map[string]interface{}{
		"key":       KEYS,
		"plaintext": string(plaintext),
		"block":     fmt.Sprintf("%d", block.BlockSize()),
		"ciphertext": fmt.Sprintf("%x", ciphertext),
		"mode":      "DES",
	}

	return fmt.Sprintf("%x", ciphertext), data, err
}

// ...

// DESDecrypt melakukan dekripsi DES
func DESDecrypt(encryptedString string, KEYS string) (decryptedString string, err error) {
	key, _ := hex.DecodeString(KEYS)
	enc, _ := hex.DecodeString(encryptedString)

	// Pastikan panjang kunci DES adalah 8 byte
	if len(key) != 8 {
		return "", errors.New("kunci DES harus memiliki panjang 8 byte")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Mode ECB untuk dekripsi DES
	mode := cipher.NewCBCDecrypter(block, make([]byte, des.BlockSize))

	plaintext := make([]byte, len(enc))
	mode.CryptBlocks(plaintext, enc)

	return string(plaintext), nil
}
