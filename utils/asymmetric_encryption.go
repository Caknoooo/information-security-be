package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateRSAKey() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

func EncryptRSA(data string, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", nil
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encData, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return "", err
	}

	return string(encData), nil
}

func DecryptRSA(data string, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", nil
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(data))
	if err != nil {
		return "", err
	}

	return string(decData), nil
}

/**
Access

rsaPriv, rsaPub, err := utils.GenerateRSAKey()
	if err != nil {
		log.Fatalf("error generate rsa key: %v", err)
	}

	fmt.Println("rsaPub: ", rsaPub)
	fmt.Println("rsaPriv: ", rsaPriv)

	sendData, err := utils.EncryptRSA("Hello World", rsaPub)
	if err != nil {
		log.Fatalf("error encrypt rsa: %v", err)
	}

	fmt.Println("sendData: ", sendData)

	decryptData, err := utils.DecryptRSA(sendData, rsaPriv)
	if err != nil {
		log.Fatalf("error decrypt rsa: %v", err)
	}

	fmt.Println("decryptData: ", decryptData)
**/