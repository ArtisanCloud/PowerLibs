package sign

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"sort"
)

type MD5Signer struct {
	Key string
}

func NewMD5Signer(key string) (signer *MD5Signer, err error) {

	signer = &MD5Signer{
		Key: key,
	}

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

func (signer *MD5Signer) KSortDataToMessage(data *object.StringMap) (msg string, err error) {

	// k sort
	var keys []string
	mapData := *data
	for k := range mapData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// join
	for _, k := range keys {
		//fmt.Println("key:", k, "Value:", mapData[k])
		msg = msg + k + "=" + mapData[k] + "&"
	}

	// omit last &
	msg = msg[0 : len(msg)-1]

	return msg, err
}

func (signer *MD5Signer) KSortObjectToMessage(data *object.HashMap) (msg string, err error) {

	// k sort
	var keys []string
	mapData := *data
	for k := range mapData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// join
	for _, k := range keys {
		//fmt.Println("key:", k, "Value:", mapData[k])
		msg = msg + fmt.Sprintf("%s=%v&", k, mapData[k])
	}

	// omit last &
	msg = msg[0 : len(msg)-1]

	return msg, err
}
