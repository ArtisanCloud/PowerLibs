package carbon

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"reflect"
)

type CarbonPeriod struct {
	startDatetime *carbon.Carbon
	endDatetime   *carbon.Carbon

	isDefaultInterval bool

	recurrences int
	options     int
}

func CreateCarbonPeriod() (p *CarbonPeriod) {

	startDatetime := carbon.Now()
	endDatetime := startDatetime.AddDay()
	p = &CarbonPeriod{
		&startDatetime,
		&endDatetime,
		true,
		0,
		0,
	}
	return p
}

func (period *CarbonPeriod) SetStartDate(date interface{}, inclusive interface{}) *CarbonPeriod {

	//fmt.Println("set start datetime")
	setDate(&period.startDatetime, date)
	return period
}

func (period *CarbonPeriod) SetEndDate(date interface{}, inclusive interface{}) *CarbonPeriod {

	//fmt.Println("set end datetime")
	setDate(&period.endDatetime, date)

	return period
}

func setDate(toSetDate **carbon.Carbon, date interface{}) (err error) {
	dType := reflect.TypeOf(date).String()
	//fmt.Printf("%v\r\n", dType)
	// 解析字符串
	if dType == "string" {
		parsedDate := carbon.Parse(date.(string))
		if parsedDate.Error == nil {

			*toSetDate = &parsedDate
		} else {
			panic("Invalid date string fmt.")
		}

	} else if dType == "carbon.Carbon" {
		// 直接赋值carbon指针
		ptr := date.(carbon.Carbon)
		*toSetDate = &ptr
	} else if dType == "*carbon.Carbon" {
		// 直接赋值carbon指针
		*toSetDate = date.(*carbon.Carbon)
	} else {
		// 如果不是string或者*carbon.Carbon， 抛出panic
		err = errors.New("Invalid date.")

	}

	return nil
}

func (period *CarbonPeriod) Overlaps(insideRange *CarbonPeriod) bool {

	//fmt.Printf("start is : %#v", period.startDatetime.ToDateTimeString())
	fmt.Printf("current start :%s %d\r\n", period.startDatetime.ToString(), period.calculateStart())
	fmt.Printf("current end   :%s %d\r\n", period.endDatetime.ToString(), period.calculateEnd())
	fmt.Printf("range start   :%s %d\r\n", insideRange.startDatetime.ToString(), insideRange.calculateStart())
	fmt.Printf("range end     :%s %d\r\n\n", insideRange.endDatetime.ToString(), insideRange.calculateEnd())

	return period.calculateEnd() > insideRange.calculateStart() && insideRange.calculateEnd() > period.calculateStart()
}

func (period *CarbonPeriod) calculateStart() int64 {
	return period.startDatetime.ToTimestamp()
}

func (period *CarbonPeriod) calculateEnd() int64 {
	return period.endDatetime.ToTimestamp()
}
