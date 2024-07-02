package youtube

import (
	"fmt"
	"os"
	"os/exec"
)

type Processor struct {
	OutputDir string
}

func NewProcessor(outputDir string) *Processor {
	return &Processor{
		OutputDir: outputDir,
	}
}

func (p *Processor) DownloadAudio(youtubeID string) (string, error) {
	const EXT = "mp3"
	outputFile := fmt.Sprintf("%s/%s.%s", p.OutputDir, youtubeID, EXT)

	// Return the file if it already exists
	if _, err := os.Stat(outputFile); err == nil {
		return outputFile, nil
	}

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputFile,
		youtubeID)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
