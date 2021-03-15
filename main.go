package main

import (
	"artisan-cloud.com/go-libs/str"
	"fmt"
)

func main() {

	strReverse := str.Reverse("hello world")
	fmt.Printf("%v", strReverse)

}
