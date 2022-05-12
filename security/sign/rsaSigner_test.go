package sign

import (
	"crypto"
	"testing"
)

func Test_RSA_VerifySignedMessage(t *testing.T) {

	signer, err := NewRSASigner(crypto.SHA256)
	if err != nil {
		t.Error(err)
	}

	_, err = signer.RSAEncryptor.GenerateKey(2048)
	if err != nil {
		t.Error(err)
	}

	a := "plane text"
	aBytes := []byte(a)

	msgHashSum, err := signer.Sign(aBytes)
	if err != nil {
		t.Error(err)
	}

	signature, err := signer.GenerateSignaturePKCS1v15(msgHashSum)
	if err != nil {
		t.Error(err)
	}

	err = signer.VerifySignPKCS1v15(msgHashSum, signature)
	if err != nil {
		t.Error(err)
	}

}
