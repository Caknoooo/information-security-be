package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"io"
)

const (
	// todo
	KEY = "8e71bbce7451ba2835de5aea73e4f3f96821455240823d2fd8174975b8321bfc!"
	FILE_KEY_AES = "8e71bbce7451ba2835de5aea73e4f3f9!"
	FILE_KEY_DES = "0011223344556677!"
	FILE_KEY_RC4 = "0011223344556677!"
)

// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes

func AESEncrypt(stringToEncrypt string, KEYS string) (encryptedString string, data map[string]interface{}, err error) {
	elapsedTimer := timerWithReturn("AESEncrypt")
	defer elapsedTimer()

	time.Sleep(1 * time.Second)

	key, _ := hex.DecodeString(KEYS)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	data = map[string]interface{}{
		"key":       KEY,
		"plaintext": string(plaintext),
		"block":     fmt.Sprintf("%d", block.BlockSize()),
    "aes-gcm":   fmt.Sprintf("%v", aesGCM), 
		"nonce":     hex.EncodeToString(nonce),
		"mode":      "AES",
		"elapsed":   elapsedTimer().String(),
	}

	return fmt.Sprintf("%x", ciphertext), data, err
}

func AESDecrypt(encryptedString string, KEYS string) (decryptedString string, err error) {
	defer func() {
		if r := recover(); r != nil {
			decryptedString = ""
			err = errors.New("error in decrypting")
		}
	}()

	key, _ := hex.DecodeString(KEYS)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", nil
	}

	return string(plaintext), nil
}