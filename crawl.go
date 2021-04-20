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
	url_prefix         = "http"
	invalid_url        = "Please specify a fully formed url starting with 'http' or 'https"
	no_url_provided    = "Please specify start page"
	exit_code          = 1
	index_zero         = 0
	minimum_args_count = 1
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	if len(args) < minimum_args_count {
		fmt.Println(no_url_provided)
		os.Exit(exit_code)
	}

	url := args[index_zero]

	if !strings.HasPrefix(url, url_prefix) {
		fmt.Println(invalid_url)
		os.Exit(exit_code)
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
