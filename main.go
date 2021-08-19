package main

import (
	"fmt"
	"github.com/ArtisanCloud/go-libs/object"
)

func main() {

	strCamel := object.Camel("123-ldfjl-sdf")
	fmt.Printf("%v", strCamel)

}
