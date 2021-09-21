package salt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type ReturnFunc func(string, Response) error

func readReturn(r io.Reader, fn ReturnFunc) error {
	dec := json.NewDecoder(r)

	// A well formed return looks like:
	//   { "return": [{ string(<id>): object(<value>), ... }] }
	err := readTokens(dec, []json.Token{
		json.Delim('{'),
		json.Token("return"),
		json.Delim('['),
		json.Delim('{'),
	})

	if err != nil {
		return err
	}

	// The data buffer is reused to minimize allocations.
	var data Response

	for dec.More() {
		t, err := dec.Token()

		if err != nil {
			return err
		}

		id, ok := t.(string)

		if !ok {
			return fmt.Errorf("expected string, received %s", t)
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

	return readTokens(dec, []json.Token{
		json.Delim('}'),
		json.Delim(']'),
		json.Delim('}'),
	})
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
