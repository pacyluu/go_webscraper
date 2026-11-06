package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Fetcher interface {
	Fetch(url string) (body string, err error)
}

type MockFetcher struct {
	MockData map[string]string
}

func (f *MockFetcher) Fetch(url string) (string, error) {
	if body, exists := f.MockData[url]; exists {
		return body, nil
	}
	return "", fmt.Errorf("not found")
}

type HTTPFetcher struct {
	client http.Client
}

func (f *HTTPFetcher) Fetch(url string) (string, error) {
	res, err := f.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetching %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return "", fmt.Errorf("fetching %s: unexpected status %d", url, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading body of %s: %w", url, err)
	}
	return string(body), nil
}

type Crawler struct {
	fetcher Fetcher
	set     map[string]bool
	queue   Queue[string]
	mu      sync.Mutex
}

func (crawler *Crawler) Crawl(startURL string, maxDepth int) {
	crawler.queue.Enqueue(startURL)
	crawler.set[startURL] = true

	currentDepth := 0

	for currentDepth < maxDepth {

		length := len(crawler.queue.elements)
		if length == 0 {
			break
		}

		batch := make([]string, 0, length)
		for range length {
			url, _ := crawler.queue.Dequeue()
			batch = append(batch, url)
		}

		var wg sync.WaitGroup
		for _, url := range batch {

			wg.Go(func() {
				fmt.Println("visiting:", url)

				body, err := crawler.fetcher.Fetch(url)
				if err != nil {
					fmt.Println("fetch error:", err)
					return
				}

				newUrls, err := ExtractLinks(url, body)
				if err != nil {
					fmt.Println("extract error:", err)
					return
				}
				for _, x := range newUrls {
					crawler.mu.Lock()
					if crawler.set[x] != true {
						crawler.set[x] = true
						crawler.queue.Enqueue(x)
					}
					crawler.mu.Unlock()
				}
			})

		}

		wg.Wait()
		currentDepth += 1
	}
}

func main() {
	c := &Crawler{
		fetcher: &HTTPFetcher{
			client: http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					MaxConnsPerHost:     20,
					MaxIdleConnsPerHost: 20,
				},
			},
		},
		set:   make(map[string]bool),
		queue: Queue[string]{},
	}

	startURL := "https://www.fragrantica.com/search/?query=Louis%20Vuitton"
	c.Crawl(startURL, 3)
}
