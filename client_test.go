package salt

import (
	"context"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	var url, username, password string

	var values = map[string]*string{
		"SALTAPI_URL":  &url,
		"SALTAPI_USER": &username,
		"SALTAPI_PASS": &password,
	}

	for k, v := range values {
		*v = os.Getenv(k)

		if *v == "" {
			t.Skipf("Skip: %s undefined", k)
		}
	}

	c := New(url)

	ctx := context.Background()

	t.Run("Login", func(t *testing.T) {
		if err := c.Login(ctx, username, password); err != nil {
			t.Fatal(err)
		}

		if len(c.Token) == 0 {
			t.Fatalf("Auth token was not set.")
		}
	})

	t.Run("Minions", func(t *testing.T) {
		var minions []string

		fn := func(id string, grains RawMessage) error {
			t.Logf("Minion: ID = %s, osfinger = %s", id, grains.Get("osfinger"))

			minions = append(minions, id)

			return nil
		}

		if err := c.Minions(ctx, fn); err != nil {
			t.Fatal(err)
		}

		if len(minions) == 0 {
			t.Fatalf("No minions were returned.")
		}

		t.Logf("Seen %d minions.", len(minions))
	})
}
