package salt

import (
	"context"
)

type PingReturnFunc func(string, bool) error

func (c *Client) Ping(ctx context.Context, target string, fn PingReturnFunc) error {
	cmd := Command{
		Client:   "local",
		Function: "test.ping",
		Target:   target,
	}

	return c.Run(ctx, &cmd, func(id string, data Response) error {
		var ok bool

		if err := data.Decode(&ok); err != nil {
			return err
		}

		return fn(id, ok)
	})
}
