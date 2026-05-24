package plugins

import (
	"fmt"
	"strings"
	"time"

	"4d63.com/tz"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

var timeZoneName string
var timezone *time.Location

func InitTimeZone(timezoneName string) {
	timezone, _ = tz.LoadLocation(timeZoneName)
}

func SetTimeZoneName(tzName string) {
	timeZoneName = tzName
}

func GetTimeZone() *time.Location {
	if timezone == nil {
		timezone, _ = tz.LoadLocation(timeZoneName)
	}

	return timezone
}

func GetNow() time.Time {
	return time.Now().In(GetTimeZone())
}

func GetNowPtr() *time.Time {
	now := GetNow()
	return &now
}

func GetCreatedAtEnd() *time.Time {
	t := GetNowPtr()
	if t != nil {
		dayEnd := t.Add(time.Hour * 24)
		dayEnd = dayEnd.Add(-time.Millisecond)
		return &dayEnd
	}
	return nil
}

func GetCreatedAtStart() *time.Time {
	t := GetNowPtr()
	if t != nil {
		startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return &startOfDay
	}
	return nil
}

type DateContract time.Time

type DateTimeContract time.Time

func (m *DateContract) UnmarshalJSON(p []byte) error {
	t, err := time.ParseInLocation(DateFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	), GetTimeZone())

	if err != nil {
		return err
	}

	*m = DateContract(t)

	return nil
}

func (m DateContract) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(m).In(GetTimeZone()).Format(DateFormat))
	return []byte(stamp), nil
}

func (m DateContract) GetTime() time.Time {
	return time.Time(m)
}

func (m *DateTimeContract) UnmarshalJSON(p []byte) error {
	t, err := time.ParseInLocation(DateTimeFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	), GetTimeZone())

	if err != nil {
		return err
	}

	*m = DateTimeContract(t)

	return nil
}

func (m DateTimeContract) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(m).In(GetTimeZone()).Format(DateTimeFormat))
	return []byte(stamp), nil
}

func (m DateTimeContract) GetTime() time.Time {
	return time.Time(m)
}

func ParseDateTimeFromQuery(value string) *time.Time {
	t, err := time.ParseInLocation(DateTimeFormat, value, GetTimeZone())
	if err != nil {
		return nil
	}
	return &t
}
