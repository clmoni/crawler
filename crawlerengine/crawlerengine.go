package crawlerengine

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/clmoni/crawler/linkengine"
)

const (
	url_format = "%s://%s"
)

// CrawlerEngine - takes a queue of links that will be iterated through
type CrawlerEngine struct {
	links   *chan string
	client  *http.Client
	visited *map[string]bool
}

func NewCrawlerEngine(l *chan string, c *http.Client, v *map[string]bool) *CrawlerEngine {
	return &CrawlerEngine{
		links:   l,
		client:  c,
		visited: v,
	}
}

// Crawl - loops thru the queue and visits each link fanning out onto child pages
func (ce *CrawlerEngine) Crawl(url string) {
	enqueue(*ce.links, url)
	host := getHostFromUrl(url)

	for link := range *ce.links {
		fullyFormedLink, err := createFullyFormedUrlIfRelative(host, link)
		if err != nil {
			fmt.Println(err.Error())
		}

		visitLink(*ce.client, host, fullyFormedLink, *ce.links, *ce.visited)
	}
}

// visitLink - visits, extracts the links and marks the page as visited
func visitLink(client http.Client, host string, link string, queue chan string, visited map[string]bool) {
	if _, found := visited[link]; !found {
		resp, err := client.Get(link)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer resp.Body.Close()
		childLinks, err := linkengine.GetLinks(resp.Body, host)
		if err != nil {
			fmt.Println(err.Error())
		}
		enqueueChildLinks(childLinks, queue)
		fmt.Println(link)
		visited[link] = true
	}
}

func enqueueChildLinks(childLinks []string, queue chan string) {
	for _, link := range childLinks {
		enqueue(queue, link)
	}
}

// asynchronous queue
func enqueue(queue chan string, link string) {
	go func() {
		queue <- link
	}()
}

// createFullyFormedUrlIfRelative - create fully formed url found in href attributes (they might be relative)
func createFullyFormedUrlIfRelative(host string, link string) (string, error) {
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	baseUrl, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	url = baseUrl.ResolveReference(url)
	return url.String(), nil
}

// getHostFromUrl - to create a well formed host to constrain the crawl
func getHostFromUrl(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf(url_format, u.Scheme, u.Hostname())
	return strings.TrimSpace(host)
}
