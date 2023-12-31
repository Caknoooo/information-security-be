package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
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

	sEnc := base64.StdEncoding.EncodeToString(encData)

	return sEnc, nil
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

	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(sDec))
	if err != nil {
		return "", err
	}

	return string(decData), nil
}

func ParsePrivateKeyFromPEM(keyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ParsePublicKeyFromPEM(pemData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Assert that the parsed key is an RSA public key
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA public key")
	}

	return publicKey, nil
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