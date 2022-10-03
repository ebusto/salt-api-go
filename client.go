package salt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	Client  *http.Client
	Headers map[string]string
	Server  string
	Token   string

	Events  *Events
	Jobs    *Jobs
	Keys    *Keys
	Minions *Minions
}

// New returns a new Client, initialized with the Salt API server.
func New(server string) *Client {
	c := &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		Server: server,
	}

	c.Events = &Events{c}
	c.Jobs = &Jobs{c}
	c.Keys = &Keys{c}
	c.Minions = &Minions{c}

	return c
}

type responseFunc func(*http.Response) error

func (c *Client) do(ctx context.Context, method, path string, data any, fn responseFunc) error {
	var buf bytes.Buffer

	if data != nil {
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.Server+"/"+path, &buf)

	if err != nil {
		return err
	}

	for key, val := range c.Headers {
		req.Header.Set(key, val)
	}

	if len(c.Token) > 0 {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return err
	}

	// Discard any unread bytes, and close the response body.
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	var ok = map[int]bool{
		200: true,
		202: true,
	}

	if !ok[res.StatusCode] {
		return NewError(res.StatusCode, htmlParagraph(res.Body))
	}

	if fn == nil {
		return nil
	}

	return fn(res)
}
