# Migration from ElevenLabs to Dia TTS

This document outlines the migration from ElevenLabs voice cloning to Dia TTS voice cloning technology.

## Summary of Changes

### Architecture Changes
- **Old**: ElevenLabs two-step process (clone voice → generate speech)
- **New**: Dia TTS one-step process (reference audio + text → generated speech)

### Files Modified

#### New Files
- `internal/zonos/client.go` - New Dia TTS API client

#### Modified Files
- `cmd/web/generateVoice.go` - Updated to use Dia TTS instead of ElevenLabs
- `README.md` - Updated documentation to reflect Dia TTS usage
- `.env` - Added `FAL_KEY` environment variable

### Key Implementation Differences

#### ElevenLabs Flow
1. Upload audio file to create voice model
2. Get voice ID from created model
3. Use voice ID to generate speech from text
4. Manage voice model lifecycle (creation/deletion)

#### Dia TTS Flow
1. Serve reference audio via local endpoint
2. Extract/provide reference text (placeholder implementation)
3. Call Dia TTS API with reference audio URL, reference text, and target text
4. Download generated audio directly

### Environment Variables
- **Removed**: `ELEVENLABS_API_KEY`
- **Added**: `FAL_KEY` (fal.ai API key)

### API Changes

#### ElevenLabs Client Methods
```go
// Old methods (removed)
func (c *Client) GetVoiceID(youtubeID string) (string, error)
func (c *Client) cloneVoice(youtubeID string) (string, error)
func (c *Client) TextToSpeech(voiceID, text string) ([]byte, error)
```

#### Dia TTS Client Methods
```go
// New methods
func (c *Client) VoiceClone(text, refAudioURL, refText string) ([]byte, error)
func (c *Client) GenerateRefAudioURL(audioFilePath, baseURL string) string
func (c *Client) ExtractReferenceText(audioFilePath string) (string, error)
```

### Benefits of Migration

1. **Simplified Workflow**: Single API call instead of multi-step process
2. **No Voice Management**: No need to manage voice model lifecycle
3. **Better Resource Usage**: No persistent voice models stored on remote service
4. **Modern API**: Uses latest fal.ai infrastructure

### TODO Items for Complete Migration

1. **Speech-to-Text**: Implement actual speech-to-text for reference text extraction
   - Current: Placeholder text `"[S1] This is sample reference text extracted from the audio."`
   - Needed: Real transcription of reference audio

2. **Text Formatting**: Enhance text formatting for Dia TTS
   - Current: Simple `[S1] {text}` format
   - Potential: Support for multi-speaker format `[S1] ... [S2] ...`

3. **Error Handling**: Enhance error handling for fal.ai API responses

4. **Testing**: Add unit tests for Dia TTS client

5. **Cleanup**: Remove unused ElevenLabs code after successful migration

### Testing the Migration

1. Start the application: `make run`
2. Navigate to `http://localhost:8080`
3. Input a YouTube/TikTok URL and text
4. Verify that voice cloning works with Dia TTS
5. Check that generated audio files are in WAV format (instead of MP3)

### Rollback Plan

If migration needs to be reversed:
1. Revert `cmd/web/generateVoice.go` to use ElevenLabs
2. Update environment variables back to `ELEVENLABS_API_KEY`
3. Remove `internal/zonos/` directory
4. Update README.md documentation

### Performance Considerations

- **Dia TTS**: Single API call, but may have longer processing time
- **Network**: Dia TTS needs to download reference audio from our server
- **File Format**: Output is WAV instead of MP3 (larger file size)

The migration maintains the same user experience while simplifying the backend implementation and reducing the complexity of voice model management.