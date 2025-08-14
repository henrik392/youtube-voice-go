package youtube

import (
	"fmt"
	"log"
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

	log.Printf("DownloadAudio: Starting download for URL: %s, VideoID: %s, OutputFile: %s", url, videoID, outputFile)

	// Return the file if it already exists
	if _, err := os.Stat(outputFile); err == nil {
		log.Printf("DownloadAudio: File already exists: %s", outputFile)
		return outputFile, nil
	}

	// Check if yt-dlp and ffmpeg are available
	ytDlpPath, err := exec.LookPath("yt-dlp")
	if err != nil {
		return "", fmt.Errorf("yt-dlp not found in PATH: %v", err)
	}
	log.Printf("DownloadAudio: Found yt-dlp at: %s", ytDlpPath)

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}
	log.Printf("DownloadAudio: Found ffmpeg at: %s", ffmpegPath)

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputFile,
		"--postprocessor-args", "ffmpeg:-t 180", // Limit to max 3 minutes
		url)

	log.Printf("DownloadAudio: Executing command: %s", cmd.String())

	// Capture both stdout and stderr for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("DownloadAudio: Command failed with error: %v", err)
		log.Printf("DownloadAudio: Command output: %s", string(output))
		return "", fmt.Errorf("yt-dlp failed: %v (output: %s)", err, string(output))
	}

	log.Printf("DownloadAudio: Command completed successfully")
	log.Printf("DownloadAudio: Output: %s", string(output))

	// Verify the output file was created
	if _, err := os.Stat(outputFile); err != nil {
		return "", fmt.Errorf("output file not created: %s (error: %v)", outputFile, err)
	}

	log.Printf("DownloadAudio: Successfully downloaded to: %s", outputFile)
	return outputFile, nil
}
