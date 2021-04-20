package crawlerengine

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/clmoni/crawler/linkengine"
)

const (
	url_format    = "%s://%s"
	forward_slash = "/"
	empty_string  = ""
)

// CrawlerEngine - takes a queue of links that will be iterated through
type CrawlerEngine struct {
	linksToVisit *chan string
	client       *http.Client
	visitedLinks *map[string]bool
	startUrl     string
}

func NewCrawlerEngine(l *chan string, c *http.Client, v *map[string]bool, u string) *CrawlerEngine {
	return &CrawlerEngine{
		linksToVisit: l,
		client:       c,
		visitedLinks: v,
		startUrl:     u,
	}
}

// Crawl - loops thru the queue and visits each link fanning out onto child pages
func (ce *CrawlerEngine) Crawl() {
	enqueue(*ce.linksToVisit, ce.startUrl)
	host := getHostFromUrl(ce.startUrl)
	for link := range *ce.linksToVisit {
		fullyFormedLink, err := createFullyFormedUrlIfRelative(host, link)
		if err != nil {
			fmt.Println(err.Error())
		}

		visitLink(ce, fullyFormedLink)
	}
}

// visitLink - visits, extracts the links and marks the page as visited
func visitLink(ce *CrawlerEngine, link string) {
	visited := *ce.visitedLinks
	if _, found := visited[link]; !found {
		go visitLinkAndEnqueueChildLinks(ce, link)
		go fmt.Println(link)
		visited[link] = true
	}
}

func visitLinkAndEnqueueChildLinks(ce *CrawlerEngine, link string) {
	resp, err := ce.client.Get(link)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	host := getHostFromUrl(ce.startUrl)
	childLinks := linkengine.GetLinks(resp.Body, host)

	go enqueueChildLinksWithinDomain(ce, childLinks)
}

func enqueueChildLinksWithinDomain(ce *CrawlerEngine, childLinks []string) {
	host := getHostFromUrl(ce.startUrl)
	for _, link := range childLinks {
		if isHrefWithinDomain(link, host) {
			enqueue(*ce.linksToVisit, link)
		}
	}
}

// isHrefWithinDomain - contrains the links found to the initial domain
func isHrefWithinDomain(link string, host string) bool {
	return strings.HasPrefix(link, forward_slash) || strings.HasPrefix(link, host)
}

// asynchronously queue
func enqueue(linksToVisit chan string, link string) {
	go func() {
		linksToVisit <- link
	}()
}

// createFullyFormedUrlIfRelative - create fully formed url found in href attributes (they might be relative)
func createFullyFormedUrlIfRelative(host string, link string) (string, error) {
	url, err := url.Parse(link)
	if err != nil {
		return empty_string, err
	}

	baseUrl, err := url.Parse(host)
	if err != nil {
		return empty_string, err
	}

	return baseUrl.ResolveReference(url).String(), nil
}

// getHostFromUrl - to create a well formed host to constrain the crawl
func getHostFromUrl(u string) string {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf(url_format, url.Scheme, url.Hostname())
	return strings.TrimSpace(host)
}
