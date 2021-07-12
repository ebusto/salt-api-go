package salt

import (
	"context"
	"net/http"
)

type Command struct {
	Arguments  []string `json:"arg,omitempty"`
	Client     string   `json:"client"`
	Function   string   `json:"fun"`
	Target     string   `json:"tgt"`
	TargetType string   `json:"tgt_type,omitempty"`
	Timeout    int      `json:"timeout,omitempty"`
}

func (c *Client) Run(ctx context.Context, cmd *Command, fn ReturnFunc) error {
	// The "local" client is the most common.
	if cmd.Client == "" {
		cmd.Client = "local"
	}

	return c.do(ctx, "POST", "/", cmd, func(r *http.Response) error {
		return readReturn(r.Body, fn)
	})
}
