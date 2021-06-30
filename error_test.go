package salt

import (
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	var tests = map[string]*Error{
		"[401] Unauthorized: msg": NewError(401, "msg"),
		"[401] Unauthorized":      NewError(401, ""),
	}

	for exp, err := range tests {
		if !reflect.DeepEqual(exp, err.Error()) {
			t.Fatalf("expected %s, received %s", exp, err.Error())
		}
	}
}
