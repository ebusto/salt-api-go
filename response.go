package salt

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

type Response []byte

func (m *Response) Decode(v any) error {
	return json.Unmarshal(*m, v)
}

func (m *Response) Delete(path string) gjson.Result {
	r := gjson.GetBytes(*m, path)

	*m, _ = sjson.DeleteBytes(*m, path)

	return r
}

func (m *Response) Get(path string) gjson.Result {
	return gjson.GetBytes(*m, path)
}

func (m *Response) Has(path string) bool {
	return m.Get(path).Exists()
}

func (m *Response) Result() gjson.Result {
	return gjson.ParseBytes(*m)
}

func (m *Response) String() string {
	return string(pretty.Pretty(*m))
}

func (m *Response) UnmarshalJSON(data []byte) error {
	*m = append((*m)[0:0], data...)

	return nil
}
