package crawler

import (
	"fmt"
	"net/url"

	"github.com/goware/urlx"
)

// LinkProcessor
type LinkProcessor interface {
	Process(source *url.URL, link *Link)
}

func NewLinkProcessor(domains []string) LinkProcessor {
	return &linkProcessor{
		domains: domains,
	}
}

type linkProcessor struct {
	domains []string
}

// Process resolves raw URL based on source URL.
func (processor *linkProcessor) Process(source *url.URL, link *Link) {
	sourceUrl, _ := urlx.Parse(source.String())
	// parse ref
	refUrl, err := url.Parse(link.RawRef)
	if err != nil {
		link.SetMalformed(fmt.Sprintf("Couldn't parse ref: %s", err.Error()))
		return
	}

	// to simplify, just remove fragments to skip fetch duplicated pages
	if len(refUrl.Fragment) > 0 {
		refUrl.Fragment = ""
	}

	link.Url = sourceUrl.ResolveReference(refUrl)
}
