package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateTimeFormat struct {
	time.Time
}

func (dateTimeFormat *DateTimeFormat) MarshalJSON() ([]byte, error) {
	formatted := dateTimeFormat.Format(time.DateTime)
	return []byte(`"` + formatted + `"`), nil
}

func (dateTimeFormat *DateTimeFormat) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(time.DateTime, string(b))
	if err != nil {
		return err
	}
	dateTimeFormat.Time = parsed
	return nil
}

func (dateTimeFormat *DateTimeFormat) Value() (driver.Value, error) {
	return dateTimeFormat.Time, nil
}

func (dateTimeFormat *DateTimeFormat) Scan(value interface{}) error {
	if value == nil {
		*dateTimeFormat = DateTimeFormat{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*dateTimeFormat = DateTimeFormat{Time: v}
	case string:
		parsedTime, err := time.Parse(time.DateTime, v)
		if err != nil {
			return err
		}
		*dateTimeFormat = DateTimeFormat{Time: parsedTime}
	default:
		return fmt.Errorf("cannot scan type %T into CustomTime", value)
	}

	return nil
}

func (dateTimeFormat *DateTimeFormat) String() string {
	return dateTimeFormat.Format(time.DateTime)
}
