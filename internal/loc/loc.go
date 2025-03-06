package loc

import (
	"fmt"
	"time"
)

type LocalTime struct {
	Name         string `json:"name"`
	LocationName string `json:"loc"`

	location *time.Location
	time     *time.Time
	valid    bool
}

func Init(name, locationName string) (*LocalTime, error) {
	if locationName == "" {
		return nil, fmt.Errorf("unknown timezone")
	}
	lt := LocalTime{Name: name, LocationName: locationName}
	err := lt.LoadLocation(lt.LocationName)
	if err != nil {
		return nil, err
	}
	lt.valid = true
	return &lt, nil
}

func (lt *LocalTime) Valid() bool {
	return lt.valid
}

func (lt *LocalTime) Validate() {
	err := lt.LoadLocation(lt.LocationName)
	if err != nil {
		return
	}
	lt.valid = true
}

func (lt *LocalTime) SetTime(t *time.Time) {
	lt.time = t
}

func (lt *LocalTime) TimeString() string {
	t := lt.time.In(lt.location)
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func (lt *LocalTime) DateString() string {
	t := lt.time.In(lt.location)
	return fmt.Sprintf("%02d-%02d-%04d", t.Month(), t.Day(), t.Year())
}

func (lt *LocalTime) LoadLocation(locationName string) error {
	tLoc, err := time.LoadLocation(locationName)
	if err != nil {
		return err
	}
	lt.location = tLoc
	return nil
}
