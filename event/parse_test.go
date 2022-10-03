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

	p.OnEvent = func(e any) error {
		fmt.Printf("event: %s\n", pretty.Sprint(e))

		return nil
	}

	p.OnUnknown = func(e salt.Response) error {
		fmt.Printf("unknown: %s\n", e.Get("tag").String())

		return nil
	}

	fh, err := os.Open("sample.json")

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

		if err := p.Parse(r); err != nil {
			t.Fatal(err)
		}
	}
}
