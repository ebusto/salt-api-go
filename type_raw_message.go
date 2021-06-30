package salt

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

type RawMessage []byte

func (m *RawMessage) Decode(v interface{}) error {
	return json.Unmarshal(*m, v)
}

func (m *RawMessage) Get(path string) gjson.Result {
	return gjson.GetBytes(*m, path)
}

func (m *RawMessage) Has(path string) bool {
	return m.Get(path).Exists()
}

func (m *RawMessage) String() string {
	return string(pretty.Pretty(*m))
}

func (m *RawMessage) UnmarshalJSON(data []byte) error {
	*m = append((*m)[0:0], data...)

	return nil
}
