package event

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

const layout = "2006-01-02T15:04:05.999999"

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var val string

	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	t.Time, err = time.Parse(layout, val)

	return err
}
