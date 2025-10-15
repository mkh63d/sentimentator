# Sentimentator: Go Media Retrieval & Audio & Transcription Service

### Consumes Download Requests and Publishes Transcripts

## 🎯 Purpose

Handles:

-   Video/audio downloading from YouTube, TikTok, and Instagram.
-   Audio extraction.
-   Speech-to-text transcription (using Whisper or Vosk).
-   Publishing transcripts for further NLP analysis.

## 🧱 Workflow

1. Consume `video.download.request` event.
2. Download video via `yt-dlp` (detect platform automatically).
3. Extract `.wav` audio using `ffmpeg`.
4. Transcribe audio to text with a speech-to-text model (Whisper/Vosk).
5. Publish a `video.transcript.ready` event containing the transcript and metadata.

## 🧩 Event Flow

```mermaid
flowchart TD
    A[C# Orchestrator] -->|video.download.request| B[Go Downloader]
    B -->|yt-dlp + ffmpeg| C[Video Platform APIs]
    B -->|Transcribe via Whisper/Vosk| D[Transcript Text]
    B -->|video.transcript.ready| E[RabbitMQ]
    E --> F[Python NLP Analyzer]
```

## 🛠️ Technology Stack

| Component        | Technology                 | Purpose                |
| ---------------- | -------------------------- | ---------------------- |
| Framework        | Gin                        | Web server             |
| Video Processing | yt-dlp + go-ytdlp + FFmpeg | Video/audio extraction |
| Storage          | Local                      | Storage integration    |
| Queue            | amqp (RabbitMQ)            | Message queue client   |
| Scraping         | Colly                      | Web scraping framework |
| Logging          | Logrus                     | Structured logging     |
| Monitoring       | Prometheus + pprof (maybe) | Performance monitoring |

## 🔄 Processing Workflow

```mermaid
sequenceDiagram
    participant Queue as RabbitMQ
    participant Go as Go Service
    participant Platform as Video Platform
    participant FFmpeg as FFmpeg
    participant Storage as File Storage
    participant Python as Python AI Service

    Queue->>Go: video.download message
    Note over Queue,Go: {job_id, video_url, platform???}

    Go->>Platform: Download video request
    Platform-->>Go: Video stream/file

    Go->>FFmpeg: Extract audio
    FFmpeg-->>Go: Audio file (WAV/MP3)

    Go->>Storage: Upload audio file
    Storage-->>Go: File URL/Path

    Go->>Queue: Publish audio.ready message
    Note over Go,Queue: {job_id, audio_path, metadata}

    Queue->>Python: Forward to AI pipeline

    Go->>Queue: Job status update???
    Note over Go,Queue: {job_id, status: "completed"}
```
