package crawler

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrEmptyURL         = errors.New("empty URL")
	ErrNotAllowedDomain = errors.New("domain not allowed")
	ErrAlreadyCrawled   = errors.New("already crawled")
)

// New .
func New(options ...Option) *Crawler {
	c := &Crawler{
		context: context.Background(),
		cfg: &Config{
			concurrency: 5,
		},
		client:           &http.Client{Timeout: time.Second},
		stat:             NewStat(),
		uniqCrawledLinks: make(map[string]struct{}),
	}

	for _, opt := range options {
		opt(c)
	}

	c.extractor = NewExtractor(c.cfg.allowedDomains)

	return c
}

// Config .
type Config struct {
	concurrency    int
	allowedDomains []string
}

type Response struct {
	StatusCode    int
	ContentType   string
	ContentLength int
	Header        http.Header
	Body          []byte
}

type Crawler struct {
	// Crawler context to manage cancellation of crawling process.
	context context.Context
	// cfg holds Crawler configuration parameters
	cfg *Config
	// client used to fetch pages
	client *http.Client
	// limits of fetchers count
	fetchersLimit chan struct{}
	// collect Crawler statistics
	stat Stat
	// extractor
	extractor Extractor
	// run handler then new content loaded
	onFetchedHandler []func(request *http.Request, response *Response)
	// in-mem crawled links holder
	uniqCrawledLinks map[string]struct{}
	//
	linksMux sync.Mutex
	// wg holds all fetchers goroutines
	wg sync.WaitGroup
}

// Run runs crawler from startRawURL.
func (c *Crawler) Run(startRawURL string) error {
	return c.fetch(c.context, startRawURL, http.MethodGet)
}

// Visit add url to crawler queue to crawl.
func (c *Crawler) Visit(url string) error {
	return c.fetch(c.context, url, http.MethodGet)
}

// Wait waits while all fetchers exit.
func (c *Crawler) Wait() {
	c.wg.Wait()
	close(c.fetchersLimit)
}

// Stat returns crawler statistic.
func (c *Crawler) Stat() PublicStat {
	return c.stat
}

func (c *Crawler) OnFetched(handler func(request *http.Request, response *Response)) {
	c.onFetchedHandler = append(c.onFetchedHandler, handler)
}

func (c *Crawler) Extractor() Extractor {
	return c.extractor
}

func (c *Crawler) fetch(ctx context.Context, rawURL string, method string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	c.stat.AddTotalDiscovered()

	if err := c.shouldBeProcessed(rawURL, parsedURL); err != nil {
		return err
	}

	c.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
			if c.context.Err() == nil {
				c.fetchersLimit <- struct{}{}
			}
		}()

		select {
		// get fetcher from pool
		case <-c.fetchersLimit:
		case <-c.context.Done():
			return
		}

		if err := c.fetchResource(ctx, rawURL, method); err != nil {
			log.Println(err)
			return
		}
		c.stat.AddTotalFetched()
	}(&c.wg)

	return nil
}

func (c *Crawler) fetchResource(ctx context.Context, url string, method string) error {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
	}()

	// skip if status != OK
	// to simplify we skip other codes like 201, 302 etc.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url: %s, status: %d", request.URL, resp.StatusCode)
	}

	var b bytes.Buffer
	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		b.Grow(int(contentLength))
	}

	if _, err := io.Copy(&b, resp.Body); err != nil {
		return err
	}

	contentType := resp.Header.Get("Content-Type")
	// simple web crawler process HTML pages only
	if !strings.Contains(contentType, "text/html") {
		return nil
	}

	c.onFetched(request, &Response{
		StatusCode:    resp.StatusCode,
		ContentType:   contentType,
		ContentLength: b.Len(),
		Header:        resp.Header,
		Body:          b.Bytes(),
	})

	return nil
}

func (c *Crawler) onFetched(r *http.Request, resp *Response) {
	for _, handler := range c.onFetchedHandler {
		handler(r, resp)
	}
}

func (c *Crawler) shouldBeProcessed(rawURL string, url *url.URL) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	if len(c.cfg.allowedDomains) > 0 && !c.domainAllowed(url.Hostname()) {
		return fmt.Errorf("check domain '%s': %w", url.Hostname(), ErrNotAllowedDomain)
	}

	md5v := md5.Sum([]byte(rawURL))
	md5s := hex.EncodeToString(md5v[:])

	c.linksMux.Lock()
	defer c.linksMux.Unlock()
	if _, ok := c.uniqCrawledLinks[md5s]; ok {
		return fmt.Errorf("url '%s': %w", rawURL, ErrAlreadyCrawled)
	}
	c.uniqCrawledLinks[md5s] = struct{}{}
	c.stat.AddUniqDiscovered()

	// others checks

	return nil
}

func (c *Crawler) domainAllowed(domain string) bool {
	for _, allowedDomain := range c.cfg.allowedDomains {
		if domain == allowedDomain {
			return true
		}
	}

	return false
}
