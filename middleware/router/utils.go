package router

import "strings"

func parse(match string) (method string, host string, urls []string) {
	matches := strings.Split(match, " ")
	url := matches[0]
	if len(matches) > 1 {
		method = matches[0]
		url = matches[1]
	}
	urls = strings.Split(url, "/")
	host = urls[0]
	return
}
