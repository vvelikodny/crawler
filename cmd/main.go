package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/goware/urlx"

	"github.com/vvelikodny/crawler/crawler"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: crawler <start-url>")
		os.Exit(1)
	}

	startURL := args[1]
	url, err := urlx.Parse(startURL)
	if err != nil {
		fmt.Println(fmt.Errorf("parse URL '%s': %w", startURL, err))
		os.Exit(1)
	}
	fmt.Println(">>", url.Hostname())

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := crawler.New(
		crawler.WithContext(ctx),
		crawler.WithClient(client),
		crawler.WithConcurrency(5),
		crawler.WithAllowedDomains(url.Hostname()),
	)

	// use mutex to print pages and pages links in order
	var mux sync.Mutex
	c.OnFetched(func(request *http.Request, response *crawler.Response) {
		mux.Lock()
		defer mux.Unlock()
		fmt.Printf("page visited: %s\n", request.URL)

		links := c.Extractor().ExtractLinks(request.URL, response)

		for _, link := range links {
			fmt.Printf("discovered on page: %s\n", link.Url.String())
		}

		for _, link := range links {
			if !link.IsRejected() {
				c.Visit(link.Url.String())
			}
		}
	})

	if err := c.Run(url.String()); err != nil {
		fmt.Println(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(chan<- os.Signal) {
		<-sigs
		cancel()
	}(sigs)

	c.Wait()

	stat := c.Stat()
	fmt.Printf("Total: %d, Uniq: %d, Fecthed: %d\n", stat.TotalDiscovered(), stat.UniqDiscovered(), stat.TotalFetched())
}
