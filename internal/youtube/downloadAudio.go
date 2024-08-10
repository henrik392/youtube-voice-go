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

func (p *Processor) DownloadAudio(url, videoID string) (string, error) {
	const EXT = "mp3"
	outputFile := fmt.Sprintf("%s/%s.%s", p.OutputDir, videoID, EXT)

	// Return the file if it already exists
	if _, err := os.Stat(outputFile); err == nil {
		return outputFile, nil
	}

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputFile,
		"--postprocessor-args", "ffmpeg:-t 180", // Limit to max 3 minutes
		url)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
