# crawler

Simple Golang Web Crawler.

# You can
  * Crawl websites and discover links on HTML pages
  * Manage max concurrency per crawler
  * Manage allowed domains
  * See the basic crawler links statistics: total discovered, uniq discovered, fetched links

# Not implemented
  * logging with levels
  * cli to run crawler with different params like "-c 10" (concurrency)
  * additional crawler events
  * robots.txt support 
  * cookie management
  * etc.

# Build
```sh
# make build
```

# Run demo

## Docker (Ctrl+C or Cmd+C to stop crawler)
```sh
crawler# make build-docker
crawler# make run-docker
```

## Terminal (Ctrl+C or Cmd+C to stop crawler)
```sh
crawler# make
crawler# ./bin/crawler https://velikodny.com
```

#Example

```golang
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    c := crawler.New(
        crawler.WithContext(ctx),
        crawler.WithConcurrency(5),
        crawler.WithAllowedDomains("velikodny.com"),
    )
    
    c.OnFetched(func(request *http.Request, response *crawler.Response) {
        fmt.Printf("page visited: %s\n", request.URL)
        links := c.Extractor().ExtractLinks(request.URL, response)

        for _, link := range links {
            if !link.IsRejected() {
                c.Visit(link.Url.String())
            }
        }
    }
    
    c.Run("https://velikodny.com")
    c.Wait()

    stat := c.Stat()
    fmt.Printf("Total: %d, Uniq: %d, Fecthed: %d\n", stat.TotalDiscovered(), stat.UniqDiscovered(), stat.TotalFetched())
}
```