package route

import (
	"log"
	"wikipathfinder/worker/web"
)

func FindPath(source_url, dest_url string) []string {
	if source_url == dest_url {
		return []string{}
	}

	visited := make(map[string]bool)
	comeFrom := make(map[string]string)
	queue := make([]string, 1)
	queue[0] = source_url

loop:
	for len(queue) != 0 {
		cur_url := queue[0]
		queue = queue[1:]
		// check if we've already visited this url
		if visited[cur_url] {
			log.Print("Skipped:", cur_url)
			continue
		}

		// get list of hrefs from cur_url
		urls := web.GetUrlsOnPage(cur_url)
		// check if dest_url in urls
		for _, url := range urls {
			comeFrom[url] = cur_url

			if url == dest_url {
				break loop
			}
			queue = append(queue, url)
		}
		visited[cur_url] = true
	}

	path := make([]string, 0)
	cur_url := dest_url
	for cur_url != source_url {
		path = append(path, cur_url)
		cur_url = comeFrom[cur_url]
	}

	path = append(path, source_url)
	return path
}
