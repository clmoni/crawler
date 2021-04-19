package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/clmoni/crawler/crawlerengine"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	client := createHttpClientWithoutSSLVerification()
	queue := make(chan string)
	visited := make(map[string]bool)

	ce := crawlerengine.NewCrawlerEngine(&queue, &client, &visited)
	ce.Crawl(args[0])
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
