package salt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	Client  *http.Client
	Headers map[string]string
	Server  string
	Token   string
}

// New returns a new Client, initialized with the Salt API server.
func New(server string) *Client {
	return &Client{
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
}

type responseFunc func(*http.Response) error

func (c *Client) do(ctx context.Context, method, path string, body interface{}, fn responseFunc) error {
	var buf bytes.Buffer

	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
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

	defer c.DiscardAndClose(res.Body)

	var ok = map[int]bool{
		200: true,
		202: true,
	}

	if !ok[res.StatusCode] {
		return NewError(res.StatusCode, c.Paragraph(res.Body))
	}

	if fn == nil {
		return nil
	}

	return fn(res)
}

// DiscardAndClose discards any unread bytes and closes the io.ReadCloser.
func (c *Client) DiscardAndClose(r io.ReadCloser) {
	io.Copy(io.Discard, r)
	r.Close()
}

// Paragraph returns the contents of paragraph in the body of the HTML document.
// It is best effort, and will return an empty string if there is no match. The
// response body is read in its entirety, but is not closed.
func (c *Client) Paragraph(r io.Reader) string {
	var v string

	doc, err := goquery.NewDocumentFromReader(r)

	if err == nil {
		doc.Find("body p").Each(func(_ int, s *goquery.Selection) {
			v = s.Text()
		})
	}

	return v
}

// Tokens reads the expected sequence of JSON tokens from the decoder, returning
// an error if not all tokens were able to be read, or an unexpected token is
// encountered.
func (c *Client) Tokens(dec *json.Decoder, seq []json.Token) error {
	var err error
	var tok json.Token

	for _, exp := range seq {
		tok, err = dec.Token()

		if !reflect.DeepEqual(exp, tok) {
			return fmt.Errorf("expected %v (%T), received %v (%T), error %v",
				exp, exp, tok, tok, err)
		}
	}

	return err
}
