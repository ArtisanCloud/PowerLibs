package sign

import (
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"testing"
)

func Test_MD5_VerifySignedMessage(t *testing.T) {

	arrayData := &object.StringMap{
		"name":   "testName",
		"desc":   "testDesc",
		"data":   "2022-05-13 00:00:00",
		"amount": "123",
	}

	signer, err := NewMD5Signer("AppID+AppSecret")
	if err != nil {
		t.Error(err)
	}

	sortMsg, err := signer.KSortDataToMessage(arrayData)
	if err != nil {
		t.Error(err)
	}

	hashSign, err := signer.Sign(sortMsg)
	if err != nil {
		t.Error(err)
	}

	//postMsg := sortMsg + "&sign" + "=" + hashSign

	err = signer.Verify(sortMsg, hashSign)
	if err != nil {
		t.Error(err)
	}
}
