package crawler

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_shouldBeProcessed(t *testing.T) {
	crawler := New(
		WithAllowedDomains("velikodny.com"),
	)

	t.Run("empty URL", func(t *testing.T) {
		u, _ := url.Parse("")
		assert.Equal(t, ErrEmptyURL, crawler.shouldBeProcessed("", u))
	})
	t.Run("good URL", func(t *testing.T) {
		u, _ := url.Parse("https://velikodny.com")
		assert.NoError(t, crawler.shouldBeProcessed("https://velikodny.com", u))
	})
	t.Run("duplicated good URL", func(t *testing.T) {
		u, _ := url.Parse("https://velikodny.com/1")
		assert.NoError(t, crawler.shouldBeProcessed("https://velikodny.com/1", u))
		assert.True(t, errors.Is(crawler.shouldBeProcessed("https://velikodny.com/1", u), ErrAlreadyCrawled))
	})
	t.Run("not allowed domain", func(t *testing.T) {
		u, _ := url.Parse("https://velikodny1.com")
		assert.True(t, errors.Is(crawler.shouldBeProcessed("https://velikodny1.com", u), ErrNotAllowedDomain))
	})
}
