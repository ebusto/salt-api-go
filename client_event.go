package salt

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
)

type Events struct {
	*Client
}

type EventStreamFunc func(Response) error

func (c *Events) Fire(ctx context.Context, tag string, data any) error {
	req := Request{
		"client": "runner", "fun": "event.send", "tag": tag, "data": data,
	}

	return c.do(ctx, "POST", "/", req, nil)
}

func (c *Events) Stream(ctx context.Context, fn EventStreamFunc) error {
	var prefix = []byte("data: ")

	return c.do(ctx, "GET", "events", nil, func(r *http.Response) error {
		// Signal to the scanner to terminate when the context is done.
		go func() {
			<-ctx.Done()
			r.Body.Close()
		}()

		sc := bufio.NewScanner(r.Body)

		for sc.Scan() {
			data := sc.Bytes()

			if !bytes.HasPrefix(data, prefix) {
				continue
			}

			if err := fn(Response(data[len(prefix):])); err != nil {
				return err
			}
		}

		return sc.Err()
	})
}
