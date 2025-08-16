package youtube

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	// Setup cookies
	cookiesPath, err := p.setupCookies()
	if err != nil {
		return "", fmt.Errorf("failed to setup cookies: %v", err)
	}
	defer func() {
		if cookiesPath != "" {
			os.Remove(cookiesPath)
		}
	}()

	args := []string{
		"-x",
		"--audio-format", "mp3",
		"-o", outputFile,
		"--postprocessor-args", "ffmpeg:-t 180", // Limit to max 3 minutes
	}

	if cookiesPath != "" {
		args = append(args, "--cookies", cookiesPath)
	}

	args = append(args, url)
	cmd := exec.Command("yt-dlp", args...)

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

// setupCookies creates a temporary cookies file from environment variable or uses local file
func (p *Processor) setupCookies() (string, error) {
	// Check for cookies in environment variable first
	cookiesEnv := os.Getenv("YOUTUBE_COOKIES")
	if cookiesEnv != "" {
		log.Printf("Using cookies from environment variable")
		
		// Create temporary file
		tmpFile, err := os.CreateTemp("", "cookies_*.txt")
		if err != nil {
			return "", fmt.Errorf("failed to create temporary cookies file: %v", err)
		}
		defer tmpFile.Close()

		// Write cookies to temporary file
		if _, err := tmpFile.WriteString(cookiesEnv); err != nil {
			os.Remove(tmpFile.Name())
			return "", fmt.Errorf("failed to write cookies to temporary file: %v", err)
		}

		return tmpFile.Name(), nil
	}

	// Fallback to local cookies.txt file
	cookiesFile := "cookies.txt"
	if _, err := os.Stat(cookiesFile); err == nil {
		log.Printf("Using local cookies file: %s", cookiesFile)
		absPath, err := filepath.Abs(cookiesFile)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for cookies file: %v", err)
		}
		return absPath, nil
	}

	log.Printf("No cookies found - proceeding without authentication")
	return "", nil
}
