package tools

import (
	"time"
)
import "fmt"

var OneDay, _ = time.ParseDuration("24h")

func CalcDayNum(start *time.Time, end *time.Time) int {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	return int(e.Sub(s).Hours() / 24)
}

func DateFromTime(input *time.Time) time.Time {
	return time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, time.Local)
}

func AddDateByDayNum(input *time.Time, dayNum int) {
	d, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*dayNum))
	input.Add(d)
}

func EndDayOfTheMonth(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month()+1, 0, 0, 0, 0, 0, in.Location())
}

func ParseZoneTime(input string) (time.Time, error) {
	output, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", input)
	if err != nil {
		output, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", input)
		if err != nil {
			return time.Now(), err
		}
	}

	return output, nil
}