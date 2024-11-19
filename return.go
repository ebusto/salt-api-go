package salt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Format int

const (
	FormatBatch Format = iota
	FormatObject
	FormatRunner
)

type ReturnFunc func(string, Response) error

func readReturn(r io.Reader, fn ReturnFunc, format Format) error {
	handlers := map[Format]func(*json.Decoder, ReturnFunc) error{
		FormatBatch:  handleBatch,
		FormatObject: handleObject,
		FormatRunner: handleRunner,
	}

	tokens := []json.Token{
		json.Delim('{'),
		json.Token("return"),
		json.Delim('['),
	}

	// Read the initial opening sequence common to all formats.
	dec := json.NewDecoder(r)

	if err := readTokens(dec, tokens); err != nil {
		return err
	}

	// Fetch and call the appropriate handler function.
	if handler, ok := handlers[format]; ok {
		return handler(dec, fn)
	}

	return fmt.Errorf("unsupported format: %v", format)
}

// handleObject parses one object surround by keys:
//
//	{[{"m1": "res1", "m2": "res2"}]}
func handleObject(dec *json.Decoder, fn ReturnFunc) error {
	if err := readTokens(dec, []json.Token{json.Delim('{')}); err != nil {
		return err
	}

	if err := processInner(dec, fn); err != nil {
		return err
	}

	return readTokens(dec, []json.Token{
		json.Delim('}'),
		json.Delim(']'),
		json.Delim('}'),
	})
}

// handleBatch parses multiple objects from a batch return:
//
//	{[{"m1": "res1"}, {"m2": "res2"}]}
func handleBatch(dec *json.Decoder, fn ReturnFunc) error {
	for dec.More() {
		if err := readTokens(dec, []json.Token{json.Delim('{')}); err != nil {
			return err
		}

		if err := processInner(dec, fn); err != nil {
			return err
		}

		if err := readTokens(dec, []json.Token{json.Delim('}')}); err != nil {
			return err
		}
	}

	return readTokens(dec, []json.Token{
		json.Delim(']'),
		json.Delim('}'),
	})
}

// processInner handles the actual parsing of the inner keys.
func processInner(dec *json.Decoder, fn ReturnFunc) error {
	var data Response

	for dec.More() {
		t, err := dec.Token()

		if err != nil {
			return err
		}

		id, ok := t.(string)

		if !ok {
			return fmt.Errorf("expected string key, received %v", t)
		}

		if err := dec.Decode(&data); err != nil {
			return err
		}

		// Handle the "error" field, which is present when interacting with
		// newer versions of Salt. The default value is an empty object ("{}"),
		// and the format is unknown.
		if id == "error" {
			if len(data) > 2 {
				return errors.New(string(data))
			}

			continue
		}

		if fn != nil {
			if err := fn(id, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// handleRunner decodes the data by calling the runner specific parser.
func handleRunner(dec *json.Decoder, fn ReturnFunc) error {
	var data Response

	for dec.More() {
		if err := dec.Decode(&data); err != nil {
			return err
		}

		if fn != nil {
			if err := fn("", data); err != nil {
				return err
			}
		}
	}

	return readTokens(dec, []json.Token{
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
