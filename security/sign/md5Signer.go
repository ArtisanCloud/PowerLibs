package sign

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type MD5Signer struct {
	Key string
}

func NewMD5Signer() (signer *MD5Signer, err error) {

	signer = &MD5Signer{}

	return signer, err
}

func (signer *MD5Signer) Sign(msg string) (encodedSign string, err error) {

	combinedMsg := []byte(signer.Key + msg)
	digest := md5.Sum(combinedMsg)
	encodedSign = hex.EncodeToString(digest[:])

	return encodedSign, err
}

func (signer *MD5Signer) Verify(msg string, signature string) (err error) {

	signedMsg, err := signer.Sign(msg)

	//fmt.Dump(msg, signedMsg, signedMsg)
	if signature != signedMsg {
		return errors.New("failed to verify sign")
	}

	return err
}

func (signer *MD5Signer)KSortData()(msg string,err error ){


	return msg, err
}
