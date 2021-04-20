package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/clmoni/crawler/crawlerengine"
)

const (
	url_prefix = "http"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	url := args[0]

	if !strings.HasPrefix(url, url_prefix) {
		fmt.Println("Please specify a fully formed url starting with 'http' or 'https")
		os.Exit(1)
	}

	client := createHttpClientWithoutSSLVerification()
	queue := make(chan string)
	visited := make(map[string]bool)

	ce := crawlerengine.NewCrawlerEngine(&queue, &client, &visited, url)
	ce.Crawl()
}

func createHttpClientWithoutSSLVerification() http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return http.Client{
		Transport: transport,
	}
}
