package crawlerengine

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestVisitLinkNotVisited(t *testing.T) {
	client := createFakeClient("")
	linksToPrint := make(chan string, 1)
	visited := make(map[string]bool)
	var wg sync.WaitGroup

	ce := NewCrawlerEngine(&linksToPrint, &client, &visited, "", &wg)
	visitLink(ce, "test.com/tests", linksToPrint)
	visitedLinks := *ce.visitedLinks
	if !visitedLinks["test.com/tests"] {
		t.Error("Failed to visit link")
	}
}

func TestIsHrefWithinDomain(t *testing.T) {
	tables := []struct {
		link           string
		host           string
		expectedResult bool
	}{
		{"http://test.com/tests", "http://test.com/", true},
		{"/test.com/tests", "http://test.com/", true},
		{"http://test.com/tests", "http://not-test.com/", false},
	}

	for _, table := range tables {
		result := isHrefWithinDomain(table.link, table.host)

		if result != table.expectedResult {
			t.Errorf("isHrefWithinDomain returned %v, expected %v", result, table.expectedResult)
		}
	}
}

func TestCreateFullyFormedUrlIfRelative(t *testing.T) {
	tables := []struct {
		link           string
		host           string
		expectedResult string
	}{
		{"/tests", "http://test.com/", "http://test.com/tests"},
		{"http://test.com/tests", "http://test.com/", "http://test.com/tests"},
		{"http://www.test.com/tests", "http://www.test.com/", "http://www.test.com/tests"},
	}

	for _, table := range tables {
		result, err := createFullyFormedUrlIfRelative(table.host, table.link)

		if err != nil {
			t.Fatalf("createFullyFormedUrlIfRelative error %v", err)
		}

		if result != table.expectedResult {
			t.Errorf("createFullyFormedUrlIfRelative returned %v, expected %v", result, table.expectedResult)
		}
	}
}

func TestGetHostFromUrl(t *testing.T) {
	tables := []struct {
		url            string
		expectedResult string
	}{
		{"http://test.com/", "http://test.com"},
		{"http://test.com/page", "http://test.com"},
		{"http://www.test.com/page", "http://www.test.com"},
	}

	for _, table := range tables {
		result := getHostFromUrl(table.url)

		if result != table.expectedResult {
			t.Errorf("createFullyFormedUrlIfRelative returned %v, expected %v", result, table.expectedResult)
		}
	}
}

func createFakeClient(markUp string) http.Client {
	return *NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(markUp))),
			Header:     make(http.Header),
		}
	})
}
