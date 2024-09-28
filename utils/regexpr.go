package utils

import (
	"html"
	"regexp"
)

func RemoveHTMLTags(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

func DecodeHTMLEntities(input string) string {
	return html.UnescapeString(input)
}
