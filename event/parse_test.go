package event

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

		if e, ok := event.(*JobReturn); ok {
			if e.Output == "highstate" {
				res, err := e.HighState()

				if err != nil {
					log.Fatal(err)
				}

				for _, r := range res {
					log.Printf("[%s] res: %s", r.Duration, pretty.Sprint(r))
				}
			}
		}
	}
}
