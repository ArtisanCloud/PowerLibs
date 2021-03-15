package main

import (
	"artisan-cloud.com/go-libs/str"
	"fmt"
)

func main() {

	strCamel := str.Camel("123-ldfjl-sdf")
	fmt.Printf("%v", strCamel)

}
