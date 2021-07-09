package salt

import (
	"context"
	"net/http"
)

func (c *Client) Login(ctx context.Context, username, password string) error {
	req := Request{
		"username": username,
		"password": password,
		"eauth":    "pam",
	}

	return c.do(ctx, "POST", "login", req, func(r *http.Response) error {
		c.Token = r.Header.Get("X-Auth-Token")

		return nil
	})
}

func (c *Client) Logout(ctx context.Context) error {
	return c.do(ctx, "POST", "logout", nil, nil)
}
