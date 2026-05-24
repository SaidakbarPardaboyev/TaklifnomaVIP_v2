package requestmodels

import (
	"fmt"
	"strings"
	"time"

	"saidakbar.origin/plugins"
)

const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = "2006-01-02 15:04:05"
)

type DateRequestModel time.Time

type DateTimeRequestModel time.Time

func (m *DateRequestModel) UnmarshalJSON(p []byte) error {
	t, err := time.ParseInLocation(dateFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	), plugins.GetTimeZone())

	if err != nil {
		return err
	}

	*m = DateRequestModel(t)

	return nil
}

func (m DateRequestModel) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(m).In(plugins.GetTimeZone()).Format(dateFormat))
	return []byte(stamp), nil
}

func (m DateRequestModel) GetTime() time.Time {
	return time.Time(m)
}

func (m *DateTimeRequestModel) UnmarshalJSON(p []byte) error {
	t, err := time.ParseInLocation(dateTimeFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	), plugins.GetTimeZone())

	if err != nil {
		return err
	}

	*m = DateTimeRequestModel(t)

	return nil
}

func (m DateTimeRequestModel) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(m).In(plugins.GetTimeZone()).Format(dateTimeFormat))
	return []byte(stamp), nil
}

func (m DateTimeRequestModel) GetTime() time.Time {
	return time.Time(m)
}

func (m DateTimeRequestModel) GetTimePtr() *time.Time {
	t := time.Time(m)
	return &t
}

func parseDateTime(value *string) *time.Time {
	if value != nil {
		if t, err := time.ParseInLocation(dateTimeFormat, *value, plugins.GetTimeZone()); err == nil {
			return &t
		}
	}

	return nil
}
