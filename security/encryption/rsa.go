package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
	ErrNotRSAPublicKey     = errors.New("Key is not a valid RSA public key")
)

type RSAEncryptor struct {
	CertificateSerialNo string          // 签名证书序列号
	PublicKeyPath       string          // 签名公钥路径
	PublicKey           *rsa.PublicKey  // 签名公钥
	PrivateKeyPath      string          // 签名私钥路径，会自动读取出*rsa.PrivateKey
	PrivateKey          *rsa.PrivateKey // 签名私钥
	Hash                crypto.Hash
}

func NewRSAEncryptor(hash crypto.Hash) (encryptor *RSAEncryptor, err error) {

	encryptor = &RSAEncryptor{
		Hash: hash,
	}
	return encryptor, nil
}

func (encryptor *RSAEncryptor) Alg() string {
	return encryptor.Hash.String()
}

func (encryptor *RSAEncryptor) GenerateKey(bit int) (privateKey *rsa.PrivateKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return nil, err
	}
	encryptor.PrivateKey = privateKey
	encryptor.PublicKey = &(privateKey.PublicKey)
	return privateKey, err
}

func (encryptor *RSAEncryptor) LoadPrivateKeyByPath() (privateKey *rsa.PrivateKey, err error) {
	keyData, err := ioutil.ReadFile(encryptor.PrivateKeyPath)
	if err != nil {
		return privateKey, err
	}

	privateKey, err = encryptor.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return privateKey, err
	}

	encryptor.PrivateKey = privateKey
	return encryptor.PrivateKey, err

}

func (encryptor *RSAEncryptor) SavePrivateKeyByPath(path string) (err error) {
	privateKeyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	// Parse PEM block
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(encryptor.PrivateKey)
	publicKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(privateKeyFile, publicKeyBlock)
	if err != nil {
		return err
	}
	//err = privateKeyFile.Close()

	return err
}

func (encryptor *RSAEncryptor) LoadPublicKeyByPath() (publicKey *rsa.PublicKey, err error) {
	keyData, err := ioutil.ReadFile(encryptor.PublicKeyPath)
	if err != nil {
		return publicKey, err
	}

	publicKey, err = encryptor.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return publicKey, err
	}

	encryptor.PublicKey = publicKey
	return encryptor.PublicKey, err

}

func (encryptor *RSAEncryptor) SavePublicKeyByPath(path string) (err error) {
	publicKeyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	// Parse PEM block
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(encryptor.PublicKey)
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	err = pem.Encode(publicKeyFile, publicKeyBlock)
	if err != nil {
		return err
	}
	//err = publicKeyFile.Close()

	return err
}

func (encryptor *RSAEncryptor) ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

// Parse PEM encoded PKCS1 or PKCS8 private key protected with password
func (encryptor *RSAEncryptor) ParseRSAPrivateKeyFromPEMWithPassword(key []byte, password string) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}

	var blockDecrypted []byte
	if blockDecrypted, err = x509.DecryptPEMBlock(block, []byte(password)); err != nil {
		return nil, err
	}

	if parsedKey, err = x509.ParsePKCS1PrivateKey(blockDecrypted); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(blockDecrypted); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

func (encryptor *RSAEncryptor) ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {

		// 兼容 x509.ParsePKCS1PublicKey的key
		parsedKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err == nil && parsedKey != nil {
			return parsedKey.(*rsa.PublicKey), nil
		}

		// 尝试其他的证书格式 x509.ParseCertificate
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPublicKey
	}

	return pkey, nil
}

// Encrypt
// hash recommend use sha256.New() as hash
func (encryptor *RSAEncryptor) Encrypt(data []byte) ([]byte, error) {
	//hash := encryptor.Hash.New()
	return rsa.EncryptOAEP(encryptor.Hash.New(), rand.Reader, encryptor.PublicKey, data, nil)
	//return rsa.EncryptOAEP(hash, rand.Reader, &(encryptor.PrivateKey.PublicKey), data, nil)
}

// Decryption
// optHash recommend use crypto.SHA256 as hash
func (encryptor *RSAEncryptor) Decryption(ciphertext []byte) (plainText []byte, err error) {
	optHash := encryptor.Hash
	return encryptor.PrivateKey.Decrypt(nil, ciphertext, &rsa.OAEPOptions{Hash: optHash})
}
