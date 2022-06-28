package object

const WX_CURRENCY_UNIT float64 = 100

func ConvertToCentUnit(amount float64) int {
	return int(amount * WX_CURRENCY_UNIT)
}

func ConvertToYuanUnit(amount int) float64 {
	return float64(amount) / WX_CURRENCY_UNIT
}
