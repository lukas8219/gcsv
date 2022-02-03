package util

import (
	"regexp"
	"strings"

	"github.com/lukas8219/gcsv/storage"
)

func Parse(link string) string {
	preParse := link

	if strings.Contains(link, "https://") {
		return parseLink(link)
	}

	link, error := storage.Get(link)

	if error != nil {
		return preParse
	}

	return link
}

func parseLink(link string) string {
	pattern := "(https:\\/\\/docs.google.com\\/spreadsheets\\/d\\/)|(\\/edit#gid=[\\d]+)"
	regx := regexp.MustCompile(pattern)
	return regx.ReplaceAllString(link, "")
}
