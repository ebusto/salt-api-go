package event

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ebusto/salt-api-go"
	"github.com/kr/pretty"
)

func TestParse(t *testing.T) {
	p := NewParser()

	fh, err := os.Open("testdata/sample.json")

	if err != nil {
		t.Fatal(err)
	}

	dec := json.NewDecoder(fh)

	for {
		var r salt.Response

		if err := dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}

			t.Fatal(err)
		}

		event, err := p.Parse(r)

		if err != nil {
			t.Fatal(err)
		}

		if event == nil {
			fmt.Printf("unhandled: %s", r.Get("tag"))
		} else {
			fmt.Printf("event: %s", pretty.Sprint(event))
		}
	}
}
