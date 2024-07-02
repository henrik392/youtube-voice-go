package youtube

import (
	"net/url"
	"regexp"
)

func GetYoutubeId(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	// Check for the correct host
	if parsedURL.Host != "www.youtube.com" && parsedURL.Host != "youtube.com" && parsedURL.Host != "youtu.be" {
		return ""
	}

	// Adjusted regular expression to allow additional query parameters
	var videoIDRegex *regexp.Regexp
	if parsedURL.Host == "youtu.be" {
		videoIDRegex = regexp.MustCompile(`^/[a-zA-Z0-9_-]+$`)
	} else {
		videoIDRegex = regexp.MustCompile(`^/watch\?v=[a-zA-Z0-9_-]+(&[a-zA-Z0-9_-]+=[a-zA-Z0-9_-]+)*$`)
	}

	// For youtube.com, ensure the 'v' parameter exists in the query
	if parsedURL.Host == "www.youtube.com" || parsedURL.Host == "youtube.com" {
		query, _ := url.ParseQuery(parsedURL.RawQuery)
		_, exists := query["v"]
		if !videoIDRegex.MatchString(parsedURL.Path+"?"+parsedURL.RawQuery) || !exists {
			return ""
		}
		return query.Get("v")
	}

	// For youtu.be, just match the path
	if !videoIDRegex.MatchString(parsedURL.Path) {
		return ""
	}
	return parsedURL.Path[1:]
}
