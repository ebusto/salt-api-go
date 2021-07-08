package salt

import (
	"context"
	"net/http"
)

type Jobs struct {
	*Client
}

func (c *Jobs) All(ctx context.Context, fn ReturnFunc) error {
	return c.do(ctx, "GET", "jobs", nil, func(r *http.Response) error {
		return readReturn(r.Body, fn)
	})
}

func (c *Jobs) Filter(ctx context.Context, id string, fn ReturnFunc) error {
	return c.do(ctx, "GET", "jobs/"+id, nil, func(r *http.Response) error {
		return readReturn(r.Body, fn)
	})
}
