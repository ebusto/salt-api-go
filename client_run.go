package salt

import (
	"context"
	"net/http"
)

type Object map[string]any

type Command struct {
	Arguments  []string `json:"arg,omitempty"`
	Client     string   `json:"client"`
	Function   string   `json:"fun"`
	Keywords   Object   `json:"kwarg,omitempty"`
	Match      string   `json:"match,omitempty"`
	Target     string   `json:"tgt,omitempty"`
	TargetType string   `json:"tgt_type,omitempty"`
	Timeout    int      `json:"timeout,omitempty"`
	Batch      string   `json:"batch,omitempty"`
}

func (c *Client) Run(ctx context.Context, cmd *Command, fn ReturnFunc) error {
	format := FormatObject

	// The "local" client is the most common.
	if cmd.Client == "" {
		cmd.Client = "local"
	}
	if cmd.Batch != "" {
		cmd.Client = "local_batch"
	}

	if cmd.Client == "runner" {
		format = FormatRunner
	}
	if cmd.Client == "local_batch" {
		format = FormatBatch
	}

	return c.do(ctx, "POST", "/", cmd, func(r *http.Response) error {
		return readReturn(r.Body, fn, format)
	})
}
