# 📋 PROJECT: SPEAKING_HEARTS (Offline Ceremonial Translation System)

## 1. OVERVIEW & GOALS
- **Core Objective:** An open-source system for international weddings that translates speeches in real-time, displays cross-platform subtitles, and transmits synthetic audio, operating 100% offline.
- **Current State:** Early development phase (Simulation).
- **Tech Stack:**
  - **Backend/Core:** Go 1.24+ (heavy concurrency using goroutines and channels).
  - **Audio/Hardware:** `portaudio` (Go bindings) and `keybd_event` (USB events).
  - **Local AI:** `faster-whisper` (STT), NLLB-200 (Meta Translation), Piper / Coqui TTS.
  - **Network:** WebSockets (`gorilla/websocket`).
  - **Post-Event Tools:** Wails (Go + HTML/JS).
  - **Infrastructure:** Docker, Docker Compose.

## 2. DATA ARCHITECTURE (The 5 Layers)
All communication between layers happens strictly through Go channels.
1. **Hardware:** USB Microphones, PTT Buttons, Screens.
2. **Acquisition:** `portaudio` streams and keyboard event capture.
3. **Processing:** Worker pools for STT (Whisper) -> Translator (NLLB) -> TTS (Piper).
4. **Distribution:** Broadcast Manager via WebSockets for language routing.
5. **Storage:** Dual-writer recording in 30-second chunks with JSON metadata.

## 3. CORE DATA CONTRACTS (Structs)
Under no circumstances should these structs be altered without explicit approval:

```go
type AudioChunk struct {
    ID         string; Source string; Timestamp  time.Time; 
    Data       []float32; SampleRate int; LangIn string
}
type ProcessedText struct {
    OriginalChunkID string; SpeakerID string; OriginalLang string; 
    OriginalText string; Translations map[string]string; Timestamp time.Time
}
type Client struct {
    ID string; Conn *websocket.Conn; Type string; 
    PreferredLang string; Muted bool; LastSeen time.Time
}