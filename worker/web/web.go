package web

import (
	"io"
	"log"
	"net/http"
	"regexp"
)

func GetUrlsOnPage(url string) []string {
	// get page html
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// find all urls using re

	re, err := regexp.Compile(`href="/wiki/[^"]+"`)
	if err != nil {
		log.Fatal(err)
	}

	results := re.FindAll(respBody, 100 /*?*/)
	if results == nil {
		log.Fatal("Empty match")
	}
	// return found urls

	urls := make([]string, len(results))
	for i, urlbytes := range results {
		url := string(urlbytes)
		url = "https://en.wikipedia.org" + url[6:len(url)-1]
		urls[i] = url
	}

	return urls
}
