package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	url := flag.String("u", "", "Base URL")
	depth := flag.Int("d", 2, "Depth to crawl.")
	flag.Parse()

	if *url == "" {
		fmt.Println("Please specify base url.")
		os.Exit(0)
	}

	allowed_domains := []string{*url}

	fmt.Println(*url)
	// Instantiate default collector
	c := colly.NewCollector(
		// default user agent header
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0"),
		// limit crawling to the domain of the specified URL
		colly.AllowedDomains(allowed_domains...),
		// set MaxDepth to the specified depth
		colly.MaxDepth(*depth),
	)

	// Visit every link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("Link:", link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// JavaScript files
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		fmt.Println("Script src:", e.Attr("src"))
	})

	// Form action URLs
	c.OnHTML("form[action]", func(e *colly.HTMLElement) {
		fmt.Println("Form action:", e.Attr("action"))
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawling page:", r.URL.String())
	})

	if err := c.Visit("https://" + *url); err != nil {
		fmt.Println("Error on start of crawl: ", err)
	}
}
