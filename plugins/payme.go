package plugins

import "time"

func ConvertToPaymeAmount(amount float64) float64 {
	return amount * 100
}

func ConverterToPaymeTimeFormat(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.UnixMilli()
}

func ConverterFromPaymeTimeFormat(ms int64) time.Time {
	return time.Unix(ms/1000, (ms%1000)*1e6).Local()
}
