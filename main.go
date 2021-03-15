package main

import (
	"artisancloud.com/go-libs/str"
	"fmt"
)

func main() {

	strCamel := str.Camel("123-ldfjl-sdf")
	fmt.Printf("%v", strCamel)

}
