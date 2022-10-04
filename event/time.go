package event

import (
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var err error

	d.Duration, err = time.ParseDuration(string(data) + "ms")

	return err
}

type Time struct {
	time.Time
}

const layout = "\"2006-01-02T15:04:05.999999\""

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error

	t.Time, err = time.Parse(layout, string(data))

	return err
}
