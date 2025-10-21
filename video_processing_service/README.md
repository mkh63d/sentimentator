# Sentimentator: Video URL Processing Service (Go)

### Consumes Download Requests and Publishes Transcripts

## 🎯 Purpose

This service is responsible for processing video content in such a way to extract it's transcribe.

- Downloading video/audio from **YouTube**.
- Extracting audio from the downloaded files.
- Transcribing the audio to text using the **OpenAI Whisper API**.
- Publishing the resulting transcripts for downstream analysis by the NLP service.

## 🧱 Workflow

1.  **Consume** a `video.download.request` event from the message queue. The event contains the URL of content.
2.  **Download** the audio from the provided URL or video.
3.  **Transformation** ensure possession of audio format suitable for the transcription service (`.mp3`).
4.  **Transcribe** the audio to text by calling the **OpenAI Whisper API**.
5.  **Publish** a `video.transcript.ready` event to the message queue, containing the transcript text and associated metadata.

## 🧩 Event Flow

Here is a sequence diagram illustrating the interaction between the services:

```mermaid
sequenceDiagram
    participant RabbitMQ
    participant VideoProcessor as Video Processing Service (Go)
    participant YouTube
    participant WhisperAPI as OpenAI Whisper API

    RabbitMQ-->>VideoProcessor: Consume `video.download.request`
    VideoProcessor->>YouTube: Download content via yt-dlp
    YouTube-->>VideoProcessor: Audio stream/file
    VideoProcessor->>WhisperAPI: Request transcription
    WhisperAPI-->>VideoProcessor: Return transcript text
    VideoProcessor->>RabbitMQ: Publish `video.transcript.ready`
```

## 🛠️ Technology Stack

| Component              | Technology                  | Purpose & Justification                                                                                                            |
| ---------------------- | --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Web Framework          | Gin                         | Minimalistic and fast web framework allowing for dev-debugging features.                                                           |
| Video/Audio Processing | yt-dlp / go-ytdlp           | Tool for downloading video and audio from various platforms, including YouTube.                                                    |
| Audio Extraction       | FFmpeg                      | Used here to extract audio streams from video files if download of `.mp3` won't be available.                                      |
| Transcription          | OpenAI Whisper API          | Chosen based on research indicating better transcription quality, simpler implementation compared to alternatives like Google STT. |
| Message Queue          | RabbitMQ (via amqp library) | A mature message broker for asynchronous communication between services.                                                           |
| Logging                | `log/slog`                  | A Go to package from standard library.                                                                                             |
