package salt

import (
	"context"
	"net/http"
)

type Command struct {
	Arguments []string `json:"arg,omitempty"`
	Client    string   `json:"client"`
	Function  string   `json:"fun"`
	Target    string   `json:"tgt"`
}

func (c *Client) Run(ctx context.Context, cmd *Command, fn ReturnFunc) error {
	return c.do(ctx, "POST", "/", cmd, func(r *http.Response) error {
		return readReturn(r.Body, fn)
	})
}
