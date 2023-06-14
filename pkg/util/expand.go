package util

import (
	"regexp"
	"strings"
)

// Expand the text with brace syntax.
// Such as: /home/{good,bad} -> [/home/good, /home/bad]
func Expand(text string) (result []string) {
	reg := regexp.MustCompile(`\{.*\}`)
	if reg.MatchString(text) {
		brace := reg.FindString(text)
		braceItem := strings.TrimPrefix(brace, "{")
		braceItem = strings.TrimSuffix(braceItem, "}")
		items := strings.Split(braceItem, ",")

		for _, item := range items {
			result = append(result, strings.ReplaceAll(text, brace, item))
		}
	} else {
		result = []string{text}
	}
	return
}
