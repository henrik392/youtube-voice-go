package web

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func isValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check for the correct host
	if parsedURL.Host != "www.youtube.com" && parsedURL.Host != "youtube.com" && parsedURL.Host != "youtu.be" {
		return false
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
		return videoIDRegex.MatchString(parsedURL.Path+"?"+parsedURL.RawQuery) && exists
	}

	// For youtu.be, just match the path
	return videoIDRegex.MatchString(parsedURL.Path)
}

func ValidateURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	component := URLInput(isValidURL(url), url)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in ValidateURLHandler: %e", err)
	}
}
