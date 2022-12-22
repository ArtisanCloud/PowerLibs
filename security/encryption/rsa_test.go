package encryption

import (
	"crypto"
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RSA_GenerateKey(t *testing.T) {
	encryptor, err := NewRSAEncryptor(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	_, err = encryptor.GenerateKey(2048)
	if err != nil {
		t.Error(err)
	}
	fmt.Dump(encryptor.PrivateKey)
	fmt.Dump(encryptor.PublicKey)
}

func Test_RSA_SaveKeys(t *testing.T) {
	encryptor, err := NewRSAEncryptor(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	_, err = encryptor.GenerateKey(2048)
	if err != nil {
		t.Error(err)
	}

	privatePath := "test_private_key.pem"
	err = encryptor.SavePrivateKeyByPath(privatePath)
	if err != nil {
		t.Error(err)
	}
	publicPath := "test_public_key.pem"
	err = encryptor.SavePublicKeyByPath(publicPath)
	if err != nil {
		t.Error(err)
	}

	assert.FileExists(t, privatePath)
	assert.FileExists(t, publicPath)

}

func Test_RSA_ParsePrivateKeys(t *testing.T) {
	encryptor, err := NewRSAEncryptor(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	encryptor.PrivateKeyPath = "test_private_key.pem"
	privateKey, err := encryptor.LoadPrivateKeyByPath()
	if err != nil {
		t.Error(err)
	}
	if privateKey == nil {
		t.Error(errors.New("private key is nil"))
	}

}

func Test_RSA_ParsePublicKeys(t *testing.T) {
	encryptor, err := NewRSAEncryptor(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	encryptor.PublicKeyPath = "test_public_key.pem"
	publicKey, err := encryptor.LoadPublicKeyByPath()
	if err != nil {
		t.Error(err)
	}
	if publicKey == nil {
		t.Error(errors.New("public key is nil"))
	}

}

func Test_RSA_EncryptMessage(t *testing.T) {
	encryptor, err := NewRSAEncryptor(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	privateKey, err := encryptor.GenerateKey(2048)
	if err != nil {
		t.Error(err)
	}
	if privateKey == nil {
		t.Error(errors.New("private key is nil"))
	}

	a := "plane text"
	aBytes := []byte(a)

	digest, err := encryptor.Encrypt(aBytes)
	if err != nil {
		t.Error(err)
	}

	decodeMsg, err := encryptor.Decryption(digest)
	if err != nil {
		t.Error(err)
	}

	fmt.Dump(a, string(decodeMsg))
	assert.EqualValues(t, a, string(decodeMsg))

}
