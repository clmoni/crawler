package linkengine

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	a_tag_selector      = "a[href]"
	href_attr           = "href"
	forward_slash       = "/"
	doc_from_reader_err = "Error occurred creating document from reader"
)

// GetLinks - parses the html, finds all a tags & extracts the href attribute
func GetLinks(htmlReader io.Reader, host string) ([]string, error) {
	links := []string{}
	doc, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		log.Fatal(err)
		return nil, errors.New(doc_from_reader_err)
	}

	doc.Find(a_tag_selector).Each(func(index int, item *goquery.Selection) {
		href, found := item.Attr(href_attr)
		if found && isHrefWithinDomain(href, host) {
			links = append(links, href)
		}
	})

	return links, nil
}

// isHrefWithinDomain - contrains the links found to the initial domain
func isHrefWithinDomain(href string, host string) bool {
	return strings.HasPrefix(href, forward_slash) || strings.HasPrefix(href, host)
}
