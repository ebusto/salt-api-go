package salt

import (
	"testing"
)

var testData = []byte(`
	{ "id": "minion" }
`)

type testStruct struct {
	ID string `json:"id"`
}

func TestResponse(t *testing.T) {
	var r = new(Response)
	var s testStruct

	if err := r.UnmarshalJSON(testData); err != nil {
		t.Fatalf("UnmarshalJSON: %s", err)
	}

	if !r.Has("id") {
		t.Fatal("Has id = false")
	}

	if r.Get("id").String() != "minion" {
		t.Fatalf("Get 'id', expected minion, received %s", r.Get("id"))
	}

	if err := r.Decode(&s); err != nil {
		t.Fatalf("Decode: %s", err)
	}

	if s.ID != "minion" {
		t.Fatalf("Struct ID = %s", s.ID)
	}
}
