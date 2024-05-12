package salt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	var tests = map[string]*Error{
		"[401] Unauthorized: msg": NewError(401, "msg"),
		"[401] Unauthorized":      NewError(401, ""),
	}

	for exp, err := range tests {
		assert.Equal(t, exp, err.Error())
	}
}
