package youtube

import (
	"regexp"
)

func GetYoutubeID(url string) string {
	patterns := []string{
		`(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`,
		`(?:youtube\.com\/shorts\/)([^"&?\/\s]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}
