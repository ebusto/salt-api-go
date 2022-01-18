package salt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type Format int

const (
	FormatObject Format = iota
	FormatRunner
)

type ReturnFunc func(string, Response) error

func readReturn(r io.Reader, fn ReturnFunc, format Format) error {
	dec := json.NewDecoder(r)

	// Object return:
	//   { "return": [ { string(<id>): object(<value>), ... } ] }
	//
	// Runner return:
	//   { "return": [ <value> ] }
	var tokens = []json.Token{
		json.Delim('{'),
		json.Token("return"),
		json.Delim('['),
	}

	if format == FormatObject {
		tokens = append(tokens, json.Delim('{'))
	}

	if err := readTokens(dec, tokens); err != nil {
		return err
	}

	// The data buffer is reused to minimize allocations.
	var data Response

	for dec.More() {
		var id string

		if format == FormatObject {
			t, err := dec.Token()

			if err != nil {
				return err
			}

			v, ok := t.(string)

			if !ok {
				return fmt.Errorf("expected string, received %s", t)
			}

			id = v
		}

		if err := dec.Decode(&data); err != nil {
			return err
		}

		if fn != nil {
			if err := fn(id, data); err != nil {
				return err
			}
		}
	}

	switch format {
	case FormatObject:
		tokens = []json.Token{json.Delim('}'), json.Delim(']'), json.Delim('}')}

	case FormatRunner:
		tokens = []json.Token{json.Delim(']'), json.Delim('}')}
	}

	return readTokens(dec, tokens)
}

// readTokens reads the expected sequence of JSON tokens from the decoder,
// returning an error if not all tokens were able to be read, or an unexpected
// token is encountered.
func readTokens(dec *json.Decoder, seq []json.Token) error {
	var err error
	var tok json.Token

	for _, exp := range seq {
		tok, err = dec.Token()

		if !reflect.DeepEqual(exp, tok) {
			return fmt.Errorf("expected %v (%T), received %v (%T), error %v",
				exp, exp, tok, tok, err)
		}
	}

	return err
}
