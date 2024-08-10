package youtube

import (
	"regexp"
)

func ExtractVideoID(url string) string {
	patterns := map[string][]string{
		"YouTube": {
			`(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`,
			`(?:youtube\.com\/shorts\/)([^"&?\/\s]{11})`,
		},
		"TikTok": {
			`tiktok\.com\/(?:@[\w.-]+\/video\/|v\/)(\d+)`,
			`vm\.tiktok\.com\/(\w+)`,
		},
	}

	// Can also check platform
	for _, platformPatterns := range patterns {
		for _, pattern := range platformPatterns {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(url)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return ""
}
