package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/goware/urlx"
)

// LinkProcessor
type LinkProcessor interface {
	Process(link *Link)
}

func NewLinkProcessor(domains ...string) LinkProcessor {
	return &linkProcessor{
		domains: domains,
	}
}

type linkProcessor struct {
	domains []string
}

// Process resolves raw URL based on source URL.
func (processor *linkProcessor) Process(link *Link) {
	sourceUrl, _ := urlx.Parse(link.Source)
	// parse ref
	refUrl, err := url.Parse(link.Ref)
	if err != nil {
		link.SetMalformed(fmt.Sprintf("Couldn't parse ref: %s", err.Error()))
		return
	}

	link.Url = sourceUrl.ResolveReference(refUrl)

	if len(link.Url.Scheme) != 0 && link.Url.Scheme != "http" && link.Url.Scheme != "https" {
		link.SetMalformed(fmt.Sprintf("scheme %s not supported", link.Url.Scheme))
		return
	}

	// parse Ref and check if it well-formatted
	if _, err = urlx.Parse(link.Url.String()); err != nil {
		if _, err = urlx.Parse(link.Url.String()); err != nil {
			link.SetMalformed(fmt.Sprintf("Ref has wrong format: %s", err.Error()))
			return
		}
	}

	// to simplify, just remove fragments to skip fetch duplicated pages
	if len(link.Url.Fragment) > 0 {
		link.Url.Fragment = ""
	}

	link.Ref, _ = urlx.NormalizeString(strings.TrimRight(link.Url.String(), "# "))
	link.Url, _ = urlx.Parse(strings.TrimRight(link.Ref, "# "))
}
