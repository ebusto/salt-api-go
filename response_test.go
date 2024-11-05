package salt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []byte(`
	{ "id": "minion", "foo": "bar" }
`)

type testStruct struct {
	ID  string `json:"id"`
	Foo string `json:"foo"`
}

func TestResponse(t *testing.T) {
	var r = new(Response)
	var s testStruct

	assert.Nil(t, r.UnmarshalJSON(testData))

	assert.True(t, r.Has("id"))

	assert.Equal(t, r.Get("id").String(), "minion")

	assert.Nil(t, r.Decode(&s))

	assert.Equal(t, s.ID, "minion")

	assert.True(t, r.Has("foo"))

	v := r.Delete("foo")

	assert.False(t, r.Has("foo"))

	assert.Equal(t, v.String(), "bar")
}
