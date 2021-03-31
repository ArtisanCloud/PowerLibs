package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	empty = ""
	tab   = "\t"
)

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty, err
	}
	return buffer.String(), nil
}

func Dump(datas ...interface{}) {
	for _, data := range datas {
		dump(data)
	}
}

func dump(data interface{}) {
	var (
		strData string
		err     error
	)
	if reflect.TypeOf(data).Kind() != reflect.String {
		////fmt.Printf("dump data: %+v\r\n", data)
		strData = fmt.Sprintf("%+v", data)
		//fmt.Printf("dump convert string: %+v\r\n", strData)

	} else {
		strData = data.(string)
	}

	prettyJson, err := PrettyJson(strData)
	if err != nil {
		fmt.Printf("convert pretty fmt error:%v", err)
	}
	fmt.Printf("%+v", prettyJson)
}





func PrintSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
