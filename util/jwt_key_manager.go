package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

const (
	privateKeyPath = "./secret/private.pem"
	publicKeyPath  = "./secret/public.pem"
	keySize        = 2048
)

func EnsureRSAKeyPair() (*rsa.PrivateKey, error) {
	if _, err := os.Stat(privateKeyPath); errors.Is(err, os.ErrNotExist) {
		err := generateRSAKeyPair()
		if err != nil {
			return nil, err
		}
	}
	return loadPrivateKey()
}

func generateRSAKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(privateKeyPath), 0700)

	privFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	_ = pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	_ = pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})

	return nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadRSAPublicKey() (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile("./secret/public.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
