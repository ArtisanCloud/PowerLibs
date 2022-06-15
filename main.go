package main

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/filter"
	"time"
)

func main() {

	//strCamel := object.Camel("123-ldfjl-sdf")
	//fmt.Printf("%v", strCamel)
	useRetryFuncCode()

}

func useRetryFuncCode() {
	s := ""
	err := gout.GET("https://api.github.com").Debug(true).BindBody(&s).F().
		Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
		Func(func(c *gout.Context) error {
			fmt2.Dump("123")
			if c.Error != nil || c.Code == 209 {
				return filter.ErrRetry
			}

			return nil

		}).Do()

	fmt.Printf("err = %v\n", err)
}
