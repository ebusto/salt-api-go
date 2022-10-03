package event

import (
	"reflect"

	"github.com/ebusto/salt-api-go"
)

type (
	HandleEvent   func(any) error
	HandleUnknown func(salt.Response) error
)

// Parser represents an event parser.
type Parser struct {
	OnEvent   HandleEvent
	OnUnknown HandleUnknown
}

func NewParser() *Parser {
	return &Parser{
		OnEvent: func(_ any) error {
			return nil
		},

		OnUnknown: func(_ salt.Response) error {
			return nil
		},
	}
}

func (p *Parser) Parse(r salt.Response) error {
	var event struct {
		Data salt.Response `json:"data"`
		Tag  string        `json:"tag"`
	}

	if err := r.Decode(&event); err != nil {
		return err
	}

	for re, fn := range Types {
		if !re.MatchString(event.Tag) {
			continue
		}

		// Extract named captures.
		match := re.FindStringSubmatch(event.Tag)
		names := make(map[string]string)

		for n, name := range re.SubexpNames() {
			names[name] = match[n]
		}

		e := fn()

		if err := event.Data.Decode(&e); err != nil {
			return err
		}

		v := reflect.ValueOf(e).Elem()

		for i := 0; i < v.NumField(); i++ {
			// Retrieve the field value.
			f := v.Field(i)

			// Retrieve the field type, in order to access tags.
			t := v.Type().Field(i)

			// Use the tagged name as the key.
			key := t.Tag.Get("name")

			if value, ok := names[key]; ok && len(key) > 0 {
				f.SetString(value)
			}
		}

		return p.OnEvent(e)
	}

	return p.OnUnknown(r)
}
