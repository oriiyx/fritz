package slug

import (
	"regexp"
	"strings"
)

func CreateSlug(input string) string {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace one or more spaces with a single dash
	reg = regexp.MustCompile(`\s+`)
	slug := reg.ReplaceAllString(processedString, "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}
