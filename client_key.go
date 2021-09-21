package salt

import (
	"context"
)

type Keys struct {
	*Client
}

type MinionFunc func(string) error

func (c *Keys) Accept(ctx context.Context, match string, fn MinionFunc) error {
	cmd := Command{
		Client:   "wheel",
		Function: "key.accept",
		Match:    match,
	}

	return c.run(ctx, &cmd, "return.minions", fn)
}

func (c *Keys) Delete(ctx context.Context, match string) error {
	cmd := Command{
		Client:   "wheel",
		Function: "key.delete",
		Match:    match,
	}

	// Unlike other key commands, key.delete does not return a list of minions.
	return c.Client.Run(ctx, &cmd, nil)
}

func (c *Keys) ListAccepted(ctx context.Context, fn MinionFunc) error {
	return c.list(ctx, "accepted", "return.minions", fn)
}

func (c *Keys) ListPending(ctx context.Context, fn MinionFunc) error {
	return c.list(ctx, "unaccepted", "return.minions_pre", fn)
}

func (c *Keys) ListRejected(ctx context.Context, fn MinionFunc) error {
	return c.list(ctx, "rejected", "return.minions_rejected", fn)
}

// A full return for a key related command includes a nested return object, and
// the key for the list of minions depends on the key command.
//
// "return": [{
//   ...
//   "data": {
//     ...
//     "return": {
//       "minions|minions_pre|minions_rejected": [
//          "minion-01",
//          "minion-02",
//          ...
//       ]
//     }
//   }
// }]

func (c *Keys) list(ctx context.Context, status, key string, fn MinionFunc) error {
	cmd := Command{
		Client:   "wheel",
		Function: "key.list",
		Match:    status,
	}

	return c.run(ctx, &cmd, key, fn)
}

func (c *Keys) run(ctx context.Context, cmd *Command, key string, fn MinionFunc) error {
	return c.Client.Run(ctx, cmd, func(id string, data Response) error {
		if fn == nil {
			return nil
		}

		for _, r := range data.Get(key).Array() {
			if err := fn(r.String()); err != nil {
				return err
			}
		}

		return nil
	})
}
