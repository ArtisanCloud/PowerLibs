package carbon

import "github.com/golang-module/carbon"

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

	switch d := date.(type) {

	// 解析字符串
	case string:
		parsedDate := carbon.Parse(date.(string))
		if parsedDate.Error != nil {

			period.startDatetime = &parsedDate
		} else {
			panic("Invalid start date string format.")
		}

	// 直接赋值cabon指针
	case *carbon.Carbon:
		period.startDatetime = d

	// 如果不是string或者*carbon.Carbon， 抛出panic
	default:
		panic("Invalid start date.")

	}

	return period
}

func (period *CarbonPeriod) SetEndDate(date interface{}, inclusive interface{}) *CarbonPeriod {
	switch d := date.(type) {

	// 解析字符串
	case string:
		parsedDate := carbon.Parse(date.(string))
		if parsedDate.Error != nil {

			period.endDatetime = &parsedDate
		} else {
			panic("Invalid end date string format.")
		}

	// 直接赋值cabon指针
	case *carbon.Carbon:
		period.endDatetime = d

	// 如果不是string或者*carbon.Carbon， 抛出panic
	default:
		panic("Invalid end date.")

	}

	return period
}


func (period *CarbonPeriod) Overlaps(insideRange *CarbonPeriod) bool {

	return period.calculateEnd() > insideRange.calculateStart() && insideRange.calculateEnd() > period.calculateStart()
}


func (period *CarbonPeriod)calculateStart() int {
	return period.endDatetime.Millisecond()
}

func (period *CarbonPeriod)calculateEnd() int {
	return period.endDatetime.Millisecond()
}

