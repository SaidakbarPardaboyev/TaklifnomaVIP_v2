package plugins

import "time"

func FormatTime(t time.Time) string {
	return t.In(GetTimeZone()).Format(DateTimeFormat)
}

func FormatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := FormatTime(*t)
	return &s
}
