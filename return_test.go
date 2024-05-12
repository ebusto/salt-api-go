package salt

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Nil(t, readTokens(dec, seq))

	assert.NotNil(t, readTokens(dec, []json.Token{json.Delim('{')}))
}
