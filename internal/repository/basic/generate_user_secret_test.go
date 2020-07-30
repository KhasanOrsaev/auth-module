package basic

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGenerateUserSecret(t *testing.T) {
	id := "name"
	secret := "secret"
	hash := "3256dd7e08eb1ece7b422bded42afba2"

	assert.Equal(t, hash, GenerateUserSecret(id, secret))
}
