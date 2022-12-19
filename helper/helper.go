package helper

import (
	"crypto/sha256"
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v2/fmt"
	. "github.com/ArtisanCloud/PowerLibs/v2/object"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

/**
 * Transform Array Keys to Camel fmt
 *
 * @param  string $mx
 * @return strID
 */
func TransformArrayKeysToCamel(arrayData HashMap) HashMap {
	arrayTransformedKeys := HashMap{}
	for key, value := range arrayData {

		if reflect.TypeOf(value).Kind() == reflect.Map {
			value = TransformArrayKeysToCamel(value.(HashMap))
			arrayTransformedKeys[Camel(key)] = value
		}

	}
	return arrayTransformedKeys

}

/**
 * Transform Array Keys to Snake fmt
 *
 * @param  string $mx
 * @return strID
 */
func TransformArrayKeysToSnake(arrayData interface{}) HashMap {
	arrayTransformedKeys := HashMap{}

	for key, value := range arrayData.(HashMap) {

		if reflect.TypeOf(value).Kind() == reflect.Map {
			value = TransformArrayKeysToSnake(value.(HashMap))
			arrayTransformedKeys[Camel(key)] = value
		}

	}

	return arrayTransformedKeys

}

///*
// * This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
// * Example how to use posted in sample_test.go file.
// */
//func StructToMap(item interface{}) map[string]interface{} {
//
//	res := map[string]interface{}{}
//	if item == nil {
//		return res
//	}
//	v := reflect.TypeOf(item)
//	reflectValue := reflect.ValueOf(item)
//	reflectValue = reflect.Indirect(reflectValue)
//
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	for i := 0; i < v.NumField(); i++ {
//		tag := v.Field(i).Tag.Get("json")
//		field := reflectValue.Field(i).Interface()
//		if tag != "" && tag != "-" {
//			if v.Field(i).Type.Kind() == reflect.Struct {
//				res[tag] = StructToMap(field)
//			} else {
//				res[tag] = field
//			}
//		}
//	}
//	return res
//}

func EncodePassword(plainPassword string) (encodedPassword string, err error) {

	encoded := sha256.Sum256([]byte(plainPassword))
	encodedPassword = fmt.Sprintf("%x", encoded)

	//fmt.Println("encoded password", encodedPassword)

	return encodedPassword, err

}

func HashPassword(encodedPassword string) (hashedPassword string, err error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(encodedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	//fmt.Println("Hashed password", string(hashed))

	return string(hashed), err
}

func EncodePlainPassword(plainPassword string) (hashedPassword string, err error) {
	encodedPassword, err := EncodePassword(plainPassword)
	if err != nil {
		return encodedPassword, err
	}
	hashedPassword, err = HashPassword(encodedPassword)

	return hashedPassword, err
}

func CheckPassword(hashedPassword string, encodedPassword string) (isPasswordValid bool, err error) {

	//fmt.Printf("hashedPassword %s\r\n", hashedPassword)
	//fmt.Printf("password %s\n", password)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(encodedPassword)); err != nil {
		fmt2.Dump(err.Error())
		return false, err
	}

	return true, nil
}
