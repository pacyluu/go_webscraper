package main

import (
	"slices"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	body := `<div><a href="/perfume/Louis-Vuitton/Pacific-Chill-81423.html" rel="noopener noreferrer" title="Louis Vuitton Pacific Chill unisex 2023" class="group prefumeHbox"><div><picture><img src="https://fimgs.net/mdimg/perfume-thumbs/270x270.81423.jpg" alt="perfume Pacific Chill"></picture></div><h3>Pacific Chill</h3></a>
<a href="#top">skip</a>
<a href="mailto:a@b.com">mail</a>
<a href="/perfume/Louis-Vuitton/Pacific-Chill-81423.html#reviews">dup</a></div>`

	got, err := ExtractLinks("https://www.fragrantica.com/search/?query=Louis%20Vuitton", body)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"https://www.fragrantica.com/perfume/Louis-Vuitton/Pacific-Chill-81423.html"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCrawlWithMockFetcher(t *testing.T) {
	start := "https://www.fragrantica.com/search/"
	page := "https://www.fragrantica.com/perfume/Louis-Vuitton/Pacific-Chill-81423.html"

	c := &Crawler{
		fetcher: &MockFetcher{MockData: map[string]string{
			start: `<a href="/perfume/Louis-Vuitton/Pacific-Chill-81423.html">Pacific Chill</a>`,
			page:  `<p>no links</p>`,
		}},
		set:   make(map[string]bool),
		queue: Queue[string]{},
	}
	c.Crawl(start, 3)

	if !c.set[page] {
		t.Errorf("expected %s to be discovered; set=%v", page, c.set)
	}
}
