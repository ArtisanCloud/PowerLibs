package sign

import (
	"testing"
)

func Test_MD5_VerifySignedMessage(t *testing.T) {

	signer, err := NewMD5Signer()
	if err != nil {
		t.Error(err)
	}

	a := "plane text"

	hashSign, err := signer.Sign(a)
	if err != nil {
		t.Error(err)
	}

	err = signer.Verify(a, hashSign)
	if err != nil {
		t.Error(err)
	}
}
