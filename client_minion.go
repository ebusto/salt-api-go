package salt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MinionFunc func(id string, grains Response) error

func (c *Client) Minions(ctx context.Context, fn MinionFunc) error {
	return c.do(ctx, "GET", "minions", nil, func(r *http.Response) error {
		dec := json.NewDecoder(r.Body)

		// A well formed response looks like:
		//   { "return": [{ string(<id>): object(<grains>) }] }
		seq := []json.Token{
			json.Delim('{'),
			json.Token("return"),
			json.Delim('['),
			json.Delim('{'),
		}

		if err := c.Tokens(dec, seq); err != nil {
			return err
		}

		// The grains buffer is reused between minions to minimize allocations.
		var grains Response

		for dec.More() {
			t, err := dec.Token()

			if err != nil {
				return err
			}

			id, ok := t.(string)

			if !ok {
				return fmt.Errorf("expected string(id), received %s", t)
			}

			if err := dec.Decode(&grains); err != nil {
				return err
			}

			if err := fn(id, grains); err != nil {
				return err
			}
		}

		seq = []json.Token{
			json.Delim('}'),
			json.Delim(']'),
			json.Delim('}'),
		}

		return c.Tokens(dec, seq)
	})
}
