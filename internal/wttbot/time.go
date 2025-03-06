package wttbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseLocation(locStr string) (*time.Location, error) {
	if locStr == "" {
		return nil, fmt.Errorf("unknown time zone")
	}
	return time.LoadLocation(locStr)
}

func ParseDate(date string) (res [3]int, err error) {
	dateParts := strings.Split(date, "-")
	if len(dateParts) != 3 {
		return res, fmt.Errorf("wrong date format, must be MM-DD-YYYY")
	}

	res[0], err = strconv.Atoi(dateParts[0])
	if err != nil {
		return res, fmt.Errorf("error mounth parsing: %w", err)
	}

	res[1], err = strconv.Atoi(dateParts[1])
	if err != nil {
		return res, fmt.Errorf("error day parsing: %w", err)
	}

	res[2], err = strconv.Atoi(dateParts[2])
	if err != nil {
		return res, fmt.Errorf("error year parsing: %w", err)
	}

	return res, nil
}

func ParseTime(t string) (res [2]int, err error) {
	timeParts := strings.Split(t, ":")
	if len(timeParts) != 2 {
		return res, fmt.Errorf("wrong time format, must be HH-MM")
	}

	res[0], err = strconv.Atoi(timeParts[0])
	if err != nil {
		return res, fmt.Errorf("error hour parsing: %w", err)
	}

	res[1], err = strconv.Atoi(timeParts[1])
	if err != nil {
		return res, fmt.Errorf("error minute parsing: %w", err)
	}
	return res, nil
}

func SetTime(date [3]int, t [2]int, loc *time.Location) *time.Time {
	parsedTime := time.Date(date[2], time.Month(date[0]), date[1], t[0], t[1], 0, 0, loc)
	return &parsedTime
}

func ParseDateTime(dateStr string, timeStr string, locStr string) (*time.Time, error) {
	l, err := ParseLocation(locStr)
	if err != nil {
		return nil, err
	}
	if dateStr == "" {
		dateStr = "01-01-2000"
	}
	d, err := ParseDate(dateStr)
	if err != nil {
		return nil, err
	}
	t, err := ParseTime(timeStr)
	if err != nil {
		return nil, err
	}
	return SetTime(d, t, l), nil
}
