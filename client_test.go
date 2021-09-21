package salt

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	var server, username, password string

	var values = map[string]*string{
		"SALTAPI_URL":  &server,
		"SALTAPI_USER": &username,
		"SALTAPI_PASS": &password,
	}

	for k, v := range values {
		*v = os.Getenv(k)

		if *v == "" {
			t.Skipf("Skip: %s undefined", k)
		}
	}

	c := New(server)

	ctx := context.Background()

	t.Run("Login", func(t *testing.T) {
		if err := c.Login(ctx, username, password); err != nil {
			t.Fatal(err)
		}

		if len(c.Token) == 0 {
			t.Fatalf("Auth token was not set.")
		}
	})

	t.Run("Jobs", func(t *testing.T) {
		var jobs []string

		fn := func(id string, job Response) error {
			t.Logf("Job: ID = %s, function = %s", id, job.Get("Function"))

			jobs = append(jobs, id)

			return nil
		}

		if err := c.Jobs.All(ctx, fn); err != nil {
			t.Fatal(err)
		}

		if len(jobs) == 0 {
			t.Fatalf("No jobs were returned.")
		}

		t.Logf("Seen %d jobs.", len(jobs))
	})

	t.Run("Keys", func(t *testing.T) {
		var minions []string

		fn := func(name string) error {
			minions = append(minions, name)

			return nil
		}

		if err := c.Keys.ListAccepted(ctx, fn); err != nil {
			t.Fatal(err)
		}

		if len(minions) == 0 {
			t.Fatalf("No minion keys were returned.")
		}

		t.Logf("Seen %d minion keys: %s", len(minions), minions)
	})

	t.Run("Minions", func(t *testing.T) {
		var minions []string

		fn := func(id string, grains Response) error {
			t.Logf("Minion: ID = %s, osfinger = %s", id, grains.Get("osfinger"))

			minions = append(minions, id)

			return nil
		}

		if err := c.Minions.All(ctx, fn); err != nil {
			t.Fatal(err)
		}

		if len(minions) == 0 {
			t.Fatalf("No minions were returned.")
		}

		t.Logf("Seen %d minions.", len(minions))
	})

	t.Run("Ping", func(t *testing.T) {
		err := c.Ping(ctx, "*", func(id string, ok bool) error {
			t.Logf("Ping: ID = %s, OK = %v", id, ok)

			return nil
		})

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Events", func(t *testing.T) {
		var seen int

		ctx, cancel := context.WithTimeout(ctx, time.Second*15)

		defer cancel()

		c.Events.Stream(ctx, func(event Response) error {
			t.Logf("Events: tag = %s", event.Get("tag"))

			seen++

			return nil
		})

		if seen == 0 {
			t.Fatalf("No events were returned.")
		}
	})

	_ = time.Now()

	t.Run("Logout", func(t *testing.T) {
		if err := c.Logout(ctx); err != nil {
			t.Fatal(err)
		}

		err := c.Minions.All(ctx, func(_ string, _ Response) error {
			return nil
		})

		if err == nil {
			t.Fatalf("Expected error when listing minions after logout.")
		}
	})
}
