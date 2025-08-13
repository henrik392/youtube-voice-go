package diatts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	APIKey     string
	BaseURL    string
	S3Client   *minio.Client
	S3Bucket   string
	S3Endpoint string
}

func NewClient(apiKey string) *Client {
	// Initialize S3 client
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	s3Bucket := os.Getenv("S3_BUCKET")

	// Remove http:// or https:// from endpoint for minio client
	endpoint := strings.TrimPrefix(s3Endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	
	// Determine if we should use HTTPS based on the original endpoint
	useSSL := strings.HasPrefix(s3Endpoint, "https://")

	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKey, s3SecretKey, ""),
		Secure: useSSL, // Use HTTPS for external endpoints
	})
	if err != nil {
		log.Printf("Error creating S3 client: %v", err)
		s3Client = nil
	}

	return &Client{
		APIKey:     apiKey,
		BaseURL:    "https://fal.run/fal-ai/dia-tts/voice-clone",
		S3Client:   s3Client,
		S3Bucket:   s3Bucket,
		S3Endpoint: s3Endpoint,
	}
}

type Request struct {
	Text         string `json:"text"`
	RefAudioURL  string `json:"ref_audio_url"`
	RefText      string `json:"ref_text"`
}

type Response struct {
	Audio struct {
		URL string `json:"url"`
	} `json:"audio"`
}

func (c *Client) VoiceClone(text, refAudioFilePath, refText string) ([]byte, error) {
	log.Printf("Starting voice cloning for file: %s", refAudioFilePath)

	// First crop the audio to 15 seconds and re-encode to reduce size
	log.Printf("Cropping and compressing audio to 15 seconds...")
	croppedFilePath, err := c.cropAndCompressAudio(refAudioFilePath, 15)
	if err != nil {
		log.Printf("Error cropping audio: %v", err)
		return nil, fmt.Errorf("error cropping audio: %w", err)
	}
	defer os.Remove(croppedFilePath) // Clean up the temporary cropped file
	log.Printf("Audio cropped successfully: %s", croppedFilePath)

	// Upload cropped file to S3 and get URL
	log.Printf("Uploading audio to S3...")
	refAudioURL, err := c.uploadToS3(croppedFilePath)
	if err != nil {
		log.Printf("Error uploading audio to S3: %v", err)
		return nil, fmt.Errorf("error uploading audio to S3: %w", err)
	}
	log.Printf("S3 URL created: %s", refAudioURL)

	payload := Request{
		Text:        text,
		RefAudioURL: refAudioURL,
		RefText:     refText,
	}

	log.Printf("Marshalling request payload...")
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}
	log.Printf("Payload size: %d bytes", len(jsonPayload))

	log.Printf("Creating HTTP request to: %s", c.BaseURL)
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.APIKey))
	log.Printf("Using API key: %s...", c.APIKey[:8]) // Log only first 8 chars for security

	log.Printf("Sending request to Dia TTS API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Received response with status code: %d", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	log.Printf("Response body length: %d bytes", len(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("API request failed with status code: %d", resp.StatusCode)
		log.Printf("Response body: %s", string(body))
		return nil, fmt.Errorf("API request failed with status code: %d\nBody: %s", resp.StatusCode, body)
	}

	log.Printf("Parsing response JSON...")
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		log.Printf("Response body: %s", string(body))
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	log.Printf("Downloading generated audio from: %s", response.Audio.URL)
	return c.downloadAudio(response.Audio.URL)
}

func (c *Client) downloadAudio(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download audio: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) SaveAudioFile(audio []byte, filename string) error {
	return os.WriteFile(filename, audio, 0644)
}

func (c *Client) ExtractReferenceText(audioFilePath string) (string, error) {
	// TODO: Implement actual speech-to-text extraction
	// For now, return a placeholder that matches the Dia TTS expected format
	return "[S1] This is sample reference text extracted from the audio.", nil
}

func (c *Client) cropAndCompressAudio(inputPath string, durationSeconds int) (string, error) {
	// Create output path with _compressed suffix
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	name := base[:len(base)-len(filepath.Ext(base))]
	outputPath := filepath.Join(dir, fmt.Sprintf("%s_compressed.mp3", name))

	log.Printf("Cropping and compressing audio: %s -> %s (%d seconds)", inputPath, outputPath, durationSeconds)

	// Use ffmpeg to crop and compress the audio
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-t", fmt.Sprintf("%d", durationSeconds),
		"-acodec", "mp3",        // Ensure MP3 encoding
		"-ab", "64k",            // Lower bitrate for smaller file
		"-ar", "22050",          // Lower sample rate
		"-ac", "1",              // Mono channel
		"-y",                    // Overwrite output file
		outputPath)

	// Capture ffmpeg output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("FFmpeg command failed: %s", string(output))
		return "", fmt.Errorf("error running ffmpeg: %w (output: %s)", err, string(output))
	}

	log.Printf("Audio cropped and compressed successfully, checking file size...")

	// Check if output file was created and get its size
	if info, err := os.Stat(outputPath); err != nil {
		log.Printf("Error checking compressed file: %v", err)
		return "", fmt.Errorf("compressed file not created: %w", err)
	} else {
		log.Printf("Compressed file size: %d bytes", info.Size())
	}

	return outputPath, nil
}

func (c *Client) uploadToS3(filePath string) (string, error) {
	if c.S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// Generate unique object name
	filename := filepath.Base(filePath)
	timestamp := time.Now().Unix()
	objectName := fmt.Sprintf("audio/%d_%s", timestamp, filename)

	log.Printf("Uploading file %s to S3 bucket %s as %s", filePath, c.S3Bucket, objectName)

	// Upload file to S3
	_, err := c.S3Client.FPutObject(context.Background(), c.S3Bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: "audio/mp3",
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to S3: %w", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/%s/%s", c.S3Endpoint, c.S3Bucket, objectName)
	log.Printf("File uploaded successfully, public URL: %s", publicURL)

	return publicURL, nil
}