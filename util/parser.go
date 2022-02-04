package util

import (
	"regexp"
	"strings"
)

func ParseLink(link string) string {
	if strings.Contains(link, "https://") {
		pattern := "(https:\\/\\/docs.google.com\\/spreadsheets\\/d\\/)|(\\/edit#gid=[\\d]+)"
		regx := regexp.MustCompile(pattern)
		return regx.ReplaceAllString(link, "")
	}
	return link
}
