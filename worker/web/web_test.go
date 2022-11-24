package web

import "testing"

const TEST_URL = "https://en.wikipedia.org/wiki/Road_warrior_(computing)"

func TestGetUrlsOnPage(t *testing.T) {
	urls := GetUrlsOnPage(TEST_URL)

	for _, url := range urls {
		t.Log(url)
	}

	if len(urls) == 0 {
		t.Fail()
	}
}
