package event

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ebusto/salt-api-go"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
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

		assert.Nil(t, err)

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
