package repository_test

import (
	"github.com/stretchr/testify/assert"
	"tongla-account/ditest"
	"testing"
)

func TestCreateKeyByName(t *testing.T) {
	app := ditest.InitTestApplication(t)

	key, err := app.ApiKeysRepository.CreateKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// test create key by name 2 times
	_, err = app.ApiKeysRepository.CreateKeyByName("test-create-key")
	assert.Error(t, err)
}

func TestFindKeyByName(t *testing.T) {
	app := ditest.InitTestApplication(t)

	key, err := app.ApiKeysRepository.CreateKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	result, err := app.ApiKeysRepository.FindKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, key, app.EncryptorRepository.Decrypt(result.Secret))
}

func TestRotateKeyByName(t *testing.T) {
	app := ditest.InitTestApplication(t)

	key, err := app.ApiKeysRepository.CreateKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	newKey, err := app.ApiKeysRepository.RotateKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, newKey)

	result, err := app.ApiKeysRepository.FindKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)

	assert.NotEqual(t, key, newKey)
	assert.Equal(t, newKey, app.EncryptorRepository.Decrypt(result.Secret))
}

func TestVerifyKey(t *testing.T) {
	app := ditest.InitTestApplication(t)
	key, err := app.ApiKeysRepository.CreateKeyByName("test-create-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	result, err := app.ApiKeysRepository.VerifyKey(key)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)

	assert.Equal(t, key, app.EncryptorRepository.Decrypt(result.Secret))
	assert.Equal(t, result.Name, "test-create-key")
}
