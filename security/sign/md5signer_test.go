package sign

import (
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MD5_KSortDataToMessage(t *testing.T) {

	arrayData := &object.StringMap{
		"b": "2",
		"a": "1",
		"c": "3",
	}

	signer, err := NewMD5Signer("AppID+AppSecret")
	if err != nil {
		t.Error(err)
	}

	sortMsg, err := signer.KSortDataToMessage(arrayData)
	if err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, "a=1&b=2&c=3", sortMsg)

}

func Test_MD5_KSortObjectToMessage(t *testing.T) {

	arrayData := &object.HashMap{
		"b": 2,
		"a": "1",
		"c": 3,
	}

	signer, err := NewMD5Signer("AppID+AppSecret")
	if err != nil {
		t.Error(err)
	}

	sortMsg, err := signer.KSortObjectToMessage(arrayData)
	//fmt.Dump(sortMsg)
	if err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, "a=1&b=2&c=3", sortMsg)

}

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
