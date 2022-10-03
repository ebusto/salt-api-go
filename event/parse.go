package event

import (
	"reflect"

	"github.com/ebusto/salt-api-go"
)

// Parser represents an event parser.
type Parser struct {
	Buffer struct {
		Data salt.Response `json:"data"`
		Tag  string        `json:"tag"`
	}
}

// NewParser returns a new event parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse returns the event parsed from the Salt response. If the tag is not
// registered, then the event will be nil.
func (p *Parser) Parse(r salt.Response) (Event, error) {
	if err := r.Decode(&p.Buffer); err != nil {
		return nil, err
	}

	for re, fn := range Types {
		if !re.MatchString(p.Buffer.Tag) {
			continue
		}

		event := fn()

		// Most fields are populated when decoding into the event structure.
		if err := p.Buffer.Data.Decode(&event); err != nil {
			return nil, err
		}

		// Some fields must be populated by parsing the tag.
		match := re.FindStringSubmatch(p.Buffer.Tag)
		names := make(map[string]string)

		// Extract named captures.
		for n, name := range re.SubexpNames() {
			names[name] = match[n]
		}

		v := reflect.ValueOf(event).Elem()

		for i := 0; i < v.NumField(); i++ {
			// Retrieve the field value.
			f := v.Field(i)

			// Retrieve the field type, in order to access tags.
			t := v.Type().Field(i)

			// Use the tagged name as the key.
			k := t.Tag.Get("name")

			if v, ok := names[k]; ok && len(k) > 0 {
				f.SetString(v)
			}
		}

		return event, nil
	}

	return nil, nil
}
