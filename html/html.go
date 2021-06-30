package html

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

// Paragraph returns the contents of paragraph in the body of the HTML document.
// It is best effort, and will return an empty string if there is no match. The
// response body is read in its entirety, but is not closed.
func Paragraph(r io.Reader) string {
	var v string

	doc, err := goquery.NewDocumentFromReader(r)

	if err == nil {
		doc.Find("body p").Each(func(_ int, s *goquery.Selection) {
			v = s.Text()
		})
	}

	io.ReadAll(r)

	return v
}
