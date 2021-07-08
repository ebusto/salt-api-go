package salt

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReadTokens(t *testing.T) {
	var doc = `{ "key": [1, 2] }`

	var seq = []json.Token{
		json.Delim('{'),
		json.Token("key"),
		json.Delim('['),
		json.Token(1.0),
	}

	dec := json.NewDecoder(strings.NewReader(doc))

	if err := readTokens(dec, seq); err != nil {
		t.Fatal(err)
	}

	if err := readTokens(dec, []json.Token{json.Delim('{')}); err == nil {
		t.Fatalf("Expected error, none returned.")
	}
}
