package crawlerengine

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/clmoni/crawler/linkengine"
)

const (
	url_format          = "%s://%s"
	forward_slash       = "/"
	empty_string        = ""
	wait_group_delta    = 1
	max_concurrent_gets = 100
)

// CrawlerEngine - takes a queue of links that will be iterated through
type CrawlerEngine struct {
	linksToVisit *chan string
	linksToPrint *chan string
	client       *http.Client
	visitedLinks *map[string]bool
	startUrl     string
	waitGroup    *sync.WaitGroup
}

func NewCrawlerEngine(lp *chan string, c *http.Client, v *map[string]bool, u string, wg *sync.WaitGroup) *CrawlerEngine {
	return &CrawlerEngine{
		client:       c,
		visitedLinks: v,
		startUrl:     u,
		linksToPrint: lp,
		waitGroup:    wg,
	}
}

// Crawl - loops thru the queue and visits each link fanning out onto child pages
func (ce *CrawlerEngine) Crawl() {
	linksToVisit := make(chan string, max_concurrent_gets)
	enqueue(linksToVisit, ce.startUrl, ce.waitGroup)
	host := getHostFromUrl(ce.startUrl)

	for i := 0; i < max_concurrent_gets; i++ {
		go func() {
			for link := range linksToVisit {
				fullyFormedLink, err := createFullyFormedUrlIfRelative(host, link)
				if err != nil {
					fmt.Println("In Crawl", err.Error())
				}

				visitLink(ce, fullyFormedLink, linksToVisit)
				ce.waitGroup.Done()
			}
		}()
	}

	ce.waitGroup.Wait()
	// time.Sleep(time.Millisecond * 1000)
	// os.Exit(0)
}

// Print - print anything in the linksToPrint channel
func (ce *CrawlerEngine) Print() {
	for {
		link := <-*ce.linksToPrint
		fmt.Println(link)
	}
}

// visitLink - visits, extracts the links and marks the page as visited
func visitLink(ce *CrawlerEngine, link string, linksToVisit chan string) {
	visited := *ce.visitedLinks
	if _, found := visited[link]; !found {
		visitLinkAndEnqueueChildLinks(ce, link, linksToVisit)
		enqueue(*ce.linksToPrint, link, nil)
		visited[link] = true
	}
}

func visitLinkAndEnqueueChildLinks(ce *CrawlerEngine, link string, linksToVisit chan<- string) {
	resp, err := ce.client.Get(link)
	if err != nil {
		fmt.Println("BOMB!!!", err.Error())
	}
	defer resp.Body.Close()

	host := getHostFromUrl(ce.startUrl)
	childLinks := linkengine.GetLinks(resp.Body, host)

	enqueueChildLinksWithinDomain(ce, childLinks, linksToVisit)
}

func enqueueChildLinksWithinDomain(ce *CrawlerEngine, childLinks []string, linksToVisit chan<- string) {
	host := getHostFromUrl(ce.startUrl)
	for _, link := range childLinks {
		if isHrefWithinDomain(link, host) {
			ce.waitGroup.Add(wait_group_delta)
			enqueue(linksToVisit, link, ce.waitGroup)
		}
	}
}

// isHrefWithinDomain - contrains the links found to the initial domain
func isHrefWithinDomain(link string, host string) bool {
	return strings.HasPrefix(link, forward_slash) || strings.HasPrefix(link, host)
}

// asynchronously queue
func enqueue(c chan<- string, message string, wg *sync.WaitGroup) {
	go func() {
		if wg != nil {
			wg.Add(wait_group_delta)
		}
		c <- message
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
