package salt

import (
	"strings"
	"testing"
)

func TestHtmlParagraph(t *testing.T) {
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
		message := htmlParagraph(strings.NewReader(body))
		if expect != message {
			t.Fatalf("expected %s, received %s", expect, message)
		}
	}
}
