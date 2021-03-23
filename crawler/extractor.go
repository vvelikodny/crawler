package crawler

import (
	"bytes"
	"net/url"
	"strings"

	. "github.com/deckarep/golang-set"
	"golang.org/x/net/html"
)

var (
	resourceTags = NewSet("object", "frame", "iframe", "style", "link", "img", "script", "input", "video", "embed")
)

type Extractor interface {
	ExtractLinks(source *url.URL, content *Response) []*Link
}

func NewExtractor(domains ...string) Extractor {
	return &extractor{
		LinkProcessor: NewLinkProcessor(domains...),
	}
}

type extractor struct {
	LinkProcessor LinkProcessor
}

func (extractor *extractor) ExtractLinks(source *url.URL, response *Response) []*Link {
	results := extractor.extractLinks(source, response)

	for _, newLink := range results {
		extractor.LinkProcessor.Process(newLink)
	}

	return results
}

func (extractor *extractor) extractLinks(source *url.URL, response *Response) []*Link {
	results := make([]*Link, 0)

	z := html.NewTokenizer(bytes.NewReader(response.Body))

	var currentToken html.Token
	var linkAttrs map[string]string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return results
		case tt == html.StartTagToken, tt == html.SelfClosingTagToken:
			currentToken = z.Token()

			// parse "a" tags only
			if currentToken.Data == "a" {
				linkAttrs = extractAttrs(currentToken)
			}

			if resourceTags.Contains(currentToken.Data) {
				attrs := extractAttrs(currentToken)

				ref := ""

				if hrefVal, ok := attrs["href"]; ok {
					ref = hrefVal
				}
				if srcVal, ok := attrs["src"]; ok {
					ref = srcVal
				}

				if len(ref) > 0 {
					sourceLink, _ := NewLink(source.String())
					resourceLink := NewHrefLink(sourceLink, ref)
					results = append(results, resourceLink)
				}
			}

			// Check if the token is an <a> tag
			if currentToken.Data == "a" && len(linkAttrs["href"]) > 0 {
				sourceLink, _ := NewLink(source.String())
				link := NewHrefLink(sourceLink, linkAttrs["href"])

				results = append(results, link)
			}
		}
	}
}

func extractAttrs(token html.Token) map[string]string {
	attrs := map[string]string{}

	for _, a := range token.Attr {
		attrs[a.Key] = strings.TrimSpace(a.Val)
	}

	return attrs
}
