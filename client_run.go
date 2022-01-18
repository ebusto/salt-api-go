package salt

import (
	"context"
	"net/http"
)

type Object map[string]interface{}

type Command struct {
	Arguments  []string `json:"arg,omitempty"`
	Client     string   `json:"client"`
	Function   string   `json:"fun"`
	Keywords   Object   `json:"kwarg,omitempty"`
	Match      string   `json:"match,omitempty"`
	Target     string   `json:"tgt,omitempty"`
	TargetType string   `json:"tgt_type,omitempty"`
	Timeout    int      `json:"timeout,omitempty"`
}

func (c *Client) Run(ctx context.Context, cmd *Command, fn ReturnFunc) error {
	format := FormatObject

	// The "local" client is the most common.
	if cmd.Client == "" {
		cmd.Client = "local"
	}

	if cmd.Client == "runner" {
		format = FormatRunner
	}

	return c.do(ctx, "POST", "/", cmd, func(r *http.Response) error {
		return readReturn(r.Body, fn, format)
	})
}
