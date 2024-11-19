package salt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var server, username, password, method, target string

	var values = map[string]*string{
		"SALTAPI_URL":    &server,
		"SALTAPI_USER":   &username,
		"SALTAPI_PASS":   &password,
		"SALTAPI_EAUTH":  &method,
		"SALTAPI_TARGET": &target,
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
		assert.Nil(t, c.Login(ctx, username, password, method))

		assert.NotEmpty(t, c.Token)
	})

	t.Run("Jobs", func(t *testing.T) {
		var jobs []string

		fn := func(id string, job Response) error {
			t.Logf("Job: ID = %s, function = %s", id, job.Get("Function"))

			jobs = append(jobs, id)

			return nil
		}

		assert.Nil(t, c.Jobs.All(ctx, fn))

		assert.NotEmpty(t, jobs)

		t.Logf("Seen %d jobs.", len(jobs))
	})

	t.Run("Keys", func(t *testing.T) {
		var minions []string

		fn := func(name string) error {
			minions = append(minions, name)

			return nil
		}

		assert.Nil(t, c.Keys.ListAccepted(ctx, fn))
		assert.NotEmpty(t, minions)

		t.Logf("Seen %d minion keys: %s", len(minions), minions)
	})

	t.Run("Minions", func(t *testing.T) {
		var minions []string

		fn := func(id string, grains Response) error {
			t.Logf("Minion: ID = %s, osfinger = %s", id, grains.Get("osfinger"))

			minions = append(minions, id)

			return nil
		}

		assert.Nil(t, c.Minions.All(ctx, fn))
		assert.NotEmpty(t, minions)

		t.Logf("Seen %d minions.", len(minions))
	})

	t.Run("Ping", func(t *testing.T) {
		err := c.Ping(ctx, target, func(id string, ok bool) error {
			t.Logf("Ping: ID = %s, OK = %v", id, ok)

			return nil
		})

		assert.Nil(t, err)
	})

	t.Run("Run", func(t *testing.T) {
		type TestArgReturn struct {
			Arguments []string `json:"args"`
			Keywords  Object   `json:"kwargs"`
		}

		cmd := Command{
			Arguments: []string{"a1", "a2"},
			Function:  "test.arg_clean",
			Keywords:  Object{"k1": 1, "k2": false},
			Target:    target,
			Timeout:   10,
		}

		exp := TestArgReturn{
			Arguments: []string{"a1", "a2"},
			Keywords:  Object{"k1": float64(1), "k2": false},
		}

		err := c.Run(ctx, &cmd, func(id string, response Response) error {
			var ret TestArgReturn

			if err := response.Decode(&ret); err != nil {
				return err
			}

			t.Logf("Run: ID = %s, return = %v", id, ret)

			assert.Equal(t, exp, ret)

			return nil
		})

		assert.Nil(t, err)
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

		assert.NotEmpty(t, seen)
	})

	_ = time.Now()

	t.Run("Logout", func(t *testing.T) {
		assert.Nil(t, c.Logout(ctx))

		err := c.Minions.All(ctx, func(_ string, _ Response) error {
			return nil
		})

		assert.NotNil(t, err)
	})
}
