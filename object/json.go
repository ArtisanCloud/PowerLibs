package object

import "encoding/json"

func JsonEncode(v interface{}) (string, error) {
	buffer, err := json.Marshal(v)

	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func JsonDecode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JsonEscape(str string) (string, error) {
	b, err := json.Marshal(str)
	if err != nil {
		return "", err
	}
	return string(b[1 : len(b)-1]), err
}
