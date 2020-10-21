package util

import (
	"log"
	"net/url"
	"path"
)

// JoinURL ...
func JoinURL(baseURL string, pathPieces ...string) string {
	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	url.Path = path.Join(url.Path, path.Join(pathPieces...))
	return url.String()
}
