package crawler

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractor_ExtractLinks(t *testing.T) {
	extractor := NewExtractor("velikodny.com")

	source, _ := url.Parse("https://velikodny.com")
	links := extractor.ExtractLinks(source, &Response{
		Body: []byte(`
			<a href="https://velikodny.com/1">test link 1</a>
			bla-bla-bla
			<a href="https://velikodny.com/2"/>
			<img href="https://velikodny.com/3.jpg"/>
			<img href="https://velikodny.com/4.png"></img>
		`),
	})

	assert.Equal(t, 4, len(links))
	for i := range links {
		assert.Equal(t, links[i].Source, "https://velikodny.com")
	}
	assert.Equal(t, links[0].Url.String(), "https://velikodny.com/1")
	assert.Equal(t, links[1].Url.String(), "https://velikodny.com/2")
	assert.Equal(t, links[2].Url.String(), "https://velikodny.com/3.jpg")
	assert.Equal(t, links[3].Url.String(), "https://velikodny.com/4.png")
}
