package repository_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tongla-account/ditest"
)

func TestStorePassphrase(t *testing.T) {
	app := ditest.InitTestApplication(t)

	passphrase, err := app.EncryptorRepository.GetPassphrase()
	assert.NoError(t, err)
	assert.NotEmpty(t, passphrase.Hash)
	assert.Equal(t, 0, passphrase.Index)

	passphrase, err = app.EncryptorRepository.GetPassphrase()
	assert.NoError(t, err)
	assert.NotEmpty(t, passphrase.Hash)
	assert.Equal(t, 0, passphrase.Index)
}

func TestGeneratePassphrase(t *testing.T) {
	app := ditest.InitTestApplication(t)

	passphrase, err := app.EncryptorRepository.GeneratePassphrase(32)
	assert.NoError(t, err)
	assert.NotEmpty(t, passphrase)
	assert.Equal(t, len(passphrase), 32)

	passphrase2, err := app.EncryptorRepository.GeneratePassphrase(21)
	assert.NoError(t, err)
	assert.NotEmpty(t, passphrase2)
	assert.Equal(t, len(passphrase2), 21)
}

func TestEncryptAndDecrypt(t *testing.T) {
	app := ditest.InitTestApplication(t)
	message := "hello world"

	encryptedMessage := app.EncryptorRepository.Encrypt(message)
	assert.NotEmpty(t, encryptedMessage)

	decryptedMessage := app.EncryptorRepository.Decrypt(encryptedMessage)
	assert.Equal(t, message, decryptedMessage)
}
