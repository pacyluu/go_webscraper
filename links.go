package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// find anchor tags
func ExtractLinks(baseURL, body string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url %q: %w", baseURL, err)
	}

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parsing html: %w", err)
	}

	var links []string
	seen := make(map[string]struct{})

	for n := range doc.Descendants() {
		if n.Type != html.ElementNode || n.Data != "a" {
			continue
		}
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}
			href := strings.TrimSpace(attr.Val)
			if href == "" || strings.HasPrefix(href, "#") {
				break
			}
			u, err := base.Parse(href)
			if err != nil {
				break
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				break
			}
			u.Fragment = ""
			abs := u.String()
			if _, dup := seen[abs]; !dup {
				seen[abs] = struct{}{}
				links = append(links, abs)
			}
			break
		}
	}

	return links, nil
}
