package crawler

import (
	"context"
	"net/http"
)

type Option func(c *Crawler)

func WithContext(context context.Context) Option {
	return func(c *Crawler) {
		c.context = context
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Crawler) {
		c.client = client
	}
}

func WithConcurrency(threadsNum int) Option {
	return func(c *Crawler) {
		c.cfg.concurrency = threadsNum
		c.fetchersLimit = initChanCapacity(threadsNum)
	}
}

func WithAllowedDomains(allowedDomains ...string) Option {
	return func(c *Crawler) {
		for _, domain := range allowedDomains {
			c.cfg.allowedDomains = append(c.cfg.allowedDomains, domain)
		}
	}
}

// Statistic sets 3rd-party Stat interface implementation for Crawler.
func WithStatistic(stat Stat) Option {
	return func(c *Crawler) {
		c.stat = stat
	}
}

// Extractor sets 3rd-party Extractor interface implementation for Crawler.
func WithExtractor(extractor Extractor) Option {
	return func(c *Crawler) {
		c.extractor = extractor
	}
}

func initChanCapacity(threadsNum int) (sem chan struct{}) {
	sem = make(chan struct{}, threadsNum)
	for i := 0; i < threadsNum; i++ {
		sem <- struct{}{}
	}
	return
}
