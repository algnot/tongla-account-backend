package repository

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type EncryptorRepository interface {
	GetPassphrase() (entity.Encryptor, error)
	GeneratePassphrase(length int) (string, error)
	Encrypt(plaintext string) []byte
	Decrypt(ciphertext []byte) string
}

type encryptorRepository struct {
	db     *gorm.DB
	config config.AppConfig
}

func (r *encryptorRepository) storePassphrase() (entity.Encryptor, error) {
	var encryptor entity.Encryptor
	passphrase, err := r.GeneratePassphrase(32)
	encryptor.Hash = []byte(passphrase)

	if err != nil {
		return entity.Encryptor{}, err
	}

	result := r.db.Create(&encryptor)

	if result.Error != nil {
		return entity.Encryptor{}, result.Error
	}

	return encryptor, nil
}

func (r *encryptorRepository) GetPassphrase() (entity.Encryptor, error) {
	var encryptor entity.Encryptor

	secretIndex := os.Getenv("APP_SECRET_INDEX")
	secretHash := os.Getenv("APP_SECRET_HASH")

	if secretIndex != "" && secretHash != "" {
		isUseEnv := true
		secretIndex, err := strconv.Atoi(secretIndex)
		if err != nil {
			isUseEnv = false
		}

		if isUseEnv {
			encryptor = entity.Encryptor{
				Index: secretIndex,
				Hash:  []byte(secretHash),
			}
			return encryptor, nil
		}
	}

	result := r.db.Last(&encryptor)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			encryptor, err := r.storePassphrase()
			if err != nil {
				return entity.Encryptor{}, err
			}
			return encryptor, nil
		}
		return entity.Encryptor{}, result.Error
	}
	err := os.Setenv("APP_SECRET_INDEX", fmt.Sprintf("%d", encryptor.Index))
	if err != nil {
		fmt.Printf("Error setting env variable %s\n", err.Error())
		return encryptor, err
	}

	err = os.Setenv("APP_SECRET_HASH", string(encryptor.Hash))
	if err != nil {
		fmt.Printf("Error setting env variable %s\n", err.Error())
		return encryptor, err
	}
	return encryptor, nil
}

func (r *encryptorRepository) GeneratePassphrase(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytePassphrase := make([]byte, length)

	_, err := rand.Read(bytePassphrase)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		bytePassphrase[i] = charset[bytePassphrase[i]%byte(len(charset))]
	}

	return string(bytePassphrase), nil
}

func (r *encryptorRepository) Encrypt(plaintext string) []byte {
	secretKey, err := r.GetPassphrase()
	if err != nil {
		panic(err)
	}

	newCipher, err := aes.NewCipher(secretKey.Hash)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext
}

func (r *encryptorRepository) Decrypt(ciphertext []byte) string {
	secretKey, err := r.GetPassphrase()
	if err != nil {
		panic(err)
	}

	aesCipher, err := aes.NewCipher(secretKey.Hash)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}

func ProvideEncryptorRepository(db *gorm.DB, config config.AppConfig) EncryptorRepository {
	return &encryptorRepository{
		db:     db,
		config: config,
	}
}
