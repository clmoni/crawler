package linkengine

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

const (
	a_tag_selector = "a[href]"
	href_attr      = "href"
)

// GetLinks - parses the html, finds all a tags & extracts the href attribute
func GetLinks(htmlReader io.Reader, host string) []string {
	links := []string{}
	doc, err := goquery.NewDocumentFromReader(htmlReader)

	if err == nil {
		doc.Find(a_tag_selector).Each(func(index int, item *goquery.Selection) {
			href, found := item.Attr(href_attr)
			if found {
				links = append(links, href)
			}
		})
	}

	return links
}
