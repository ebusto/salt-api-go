package salt

import (
	"context"
	"net/http"
)

type Minions struct {
	*Client
}

func (c *Minions) All(ctx context.Context, fn ReturnFunc) error {
	return c.do(ctx, "GET", "minions", nil, func(r *http.Response) error {
		return readReturn(r.Body, fn, FormatObject)
	})
}

func (c *Minions) Filter(ctx context.Context, id string, fn ReturnFunc) error {
	return c.do(ctx, "GET", "minions/"+id, nil, func(r *http.Response) error {
		return readReturn(r.Body, fn, FormatObject)
	})
}
