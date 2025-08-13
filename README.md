# YouTube Voice Cloner

Transform YouTube and TikTok videos into custom AI-generated speech using Dia TTS voice synthesis technology.

## What does it do?

This application takes a YouTube or TikTok video URL and creates a new audio file with the same content but spoken in a different AI-generated voice. Here's how it works:

1. **Download**: Extracts audio from YouTube/TikTok videos
2. **Analyze**: Processes the original audio for voice characteristics
3. **Synthesize**: Generates new speech using Dia TTS voice cloning technology

## Features

- ✅ **Multi-platform support**: YouTube and TikTok videos
- ✅ **Voice cloning**: Creates custom voice models from source audio
- ✅ **Web interface**: Simple, responsive UI built with HTMX
- ✅ **Real-time processing**: See progress as your audio is generated
- ✅ **Audio player**: Listen to results directly in the browser
- ✅ **Docker deployment**: Ready for cloud deployment (Google Cloud Run)

## Requirements

### System Dependencies
1. **yt-dlp** - Video/audio downloader
   ```bash
   # macOS
   brew install yt-dlp
   
   # Ubuntu/Debian
   sudo apt install yt-dlp
   
   # Or via pip
   pip install yt-dlp
   ```

2. **TailwindCSS** - For building styles
   ```bash
   npm install -g tailwindcss
   ```

3. **TEMPL** - For templating in go
    ```bash
    go install github.com/a-h/templ/cmd/templ@latest
    ```

### Environment Variables
Create a `.env` file with:
```bash
PORT=8080
FAL_KEY=your_fal_ai_api_key_here
DATABASE_URL=your_postgres_connection_string
```

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/henrik392/youtube-voice-go.git
   cd youtube-voice-go
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build and run**
   ```bash
   make build
   make run
   ```

4. **Open your browser**
   Navigate to `http://localhost:8080`

## Development Commands

```bash
# Build the application (generates templates + CSS + binary)
make build

# Run the application
make run

# Start with live reload (installs air if needed)
make watch

# Run tests
make test

# Start PostgreSQL database container
make docker-run

# Stop database container
make docker-down

# Clean build artifacts
make clean
```

## Video Limitations

- **Length**: 30 seconds to 10 minutes
- **Optimal**: 1-5 minutes with clear audio
- **Format**: Supports any format that yt-dlp can process

## Architecture

```
cmd/
├── api/           # Main application entry point
└── web/           # Web handlers and templates
internal/
├── database/      # PostgreSQL integration
├── elevenlabs/    # Voice synthesis API client
├── server/        # HTTP server setup
└── youtube/       # Video processing logic
```

## Technology Stack

- **Backend**: Go with Chi router
- **Frontend**: HTML templates (templ) + HTMX + TailwindCSS
- **Database**: PostgreSQL
- **Audio Processing**: yt-dlp + ffmpeg
- **AI Voice**: Dia TTS (fal.ai) API

## Deployment

### Docker
```bash
make docker-build
docker run -p 8080:8080 yt-voice
```

### Google Cloud Run
```bash
gcloud run deploy --image=europe-north1-docker.pkg.dev/youtube-to-voice/youtube-to-voice-repo/youtube-to-voice-image:tag1
```

## How It Works

1. **URL Validation**: Checks if the provided URL is from YouTube or TikTok
2. **Audio Extraction**: Downloads and converts video to MP3 (max 3 minutes)
3. **Reference Processing**: Prepares the original audio as reference for voice cloning
4. **Text Processing**: Formats the target text for Dia TTS processing
5. **Voice Synthesis**: Uses Dia TTS to generate speech with the cloned voice in one step
6. **Delivery**: Serves the final audio file through the web interface

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

This project is for educational and personal use. Please respect content creators' rights and fal.ai's terms of service.