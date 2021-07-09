package salt

import (
	"context"
	"os"
	"strings"
	"testing"
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

	t.Run("Paragraph", func(t *testing.T) {
		// This is the slightly simplified HTML response from Salt when a login is rejected.
		var bodyHit = `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"></meta>
                    <title>401 Unauthorized</title>
                </head>
                    <body>
                        <h2>401 Unauthorized</h2>
                        <p>Could not authenticate using provided credentials</p>
                        <pre id="traceback"></pre>
                    <div id="powered_by">
                      <span>
                        Powered by <a href="http://www.cherrypy.org">CherryPy 18.6.0</a>
                      </span>
                    </div>
                    </body>
                </html>
            `

		var bodyMiss = `
                <html>
                <head>
                    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"></meta>
                    <title>401 Unauthorized</title>
                </head>
                    <body>
                        <h2>401 Unauthorized</h2>
                        <b>Could not authenticate using provided credentials</b>
                        <pre id="traceback"></pre>
                    <div id="powered_by">
                      <span>
                        Powered by <a href="http://www.cherrypy.org">CherryPy 18.6.0</a>
                      </span>
                    </div>
                    </body>
                </html>
            `

		var bodyInvalid = `invalid html`

		var tests = map[string]string{
			bodyHit:     "Could not authenticate using provided credentials",
			bodyMiss:    "",
			bodyInvalid: "",
		}

		for body, expect := range tests {
			message := c.Paragraph(strings.NewReader(body))
			if expect != message {
				t.Fatalf("expected %s, received %s", expect, message)
			}
		}
	})
}
