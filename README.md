# 🎥 Sentimentator - Review Sentiment Analyzer

## 🧠 Overview

**Sentimentator** is a distributed, event-driven system that automatically analyzes **video reviews** of smartphones from **YouTube** and produces a **synthetic sentiment summary** of opinions about device features (e.g., camera, battery, screen, performance).

## 🌍 High-Level Architecture

The system is composed of several microservices that communicate asynchronously via a message queue.

```mermaid
graph TB
    subgraph "Client Layer"
        User[👤 User]
        Browser[🌐 Web Browser]
    end

    subgraph "Frontend Layer"
        Frontend[Vue.js Frontend<br/>Port: 80<br/>• URL Submission<br/>• Results Display<br/>• State Management]
    end

    subgraph "API Gateway Layer"
        Orchestrator[C# Core Orchestrator<br/>• REST API<br/>• Job Management<br/>• Auth & Validation<br/>• OpenAPI/Swagger]
    end

    subgraph "Message Queue Layer"
        RabbitMQ[🐰 RabbitMQ<br/>Port: 5672<br/>Management: 15672<br/>• Event Routing<br/>• Async Communication]
    end

    subgraph "Processing Services"
        VideoProc[Go Video Processor<br/>• yt-dlp Download<br/>• Audio Extraction<br/>• Whisper API Integration]
        NLP[Python NLP Analyzer<br/>• Feature Detection<br/>• Sentiment Analysis<br/>• Result Aggregation]
    end

    subgraph "External APIs"
        YouTube[📺 YouTube]
        Whisper[🎙️ OpenAI Whisper API<br/>Speech-to-Text]
    end

    subgraph "Data Layer"
        Postgres[(PostgreSQL<br/>Port: 5432<br/>• Analysis Results<br/>• User Data<br/>• Job History)]
        Redis[(Redis<br/>Cache Layer<br/>)]
    end

    User --> Browser
    Browser --> Frontend
    Frontend <-->|REST API| Orchestrator
    Orchestrator <-->|Publish/Consume| RabbitMQ
    Orchestrator <-->|EF Core| Postgres
    Orchestrator <--> Redis

    RabbitMQ <-->|video.download.request| VideoProc
    VideoProc -->|Download| YouTube
    VideoProc -->|Transcribe| Whisper
    VideoProc -->|video.transcript.ready| RabbitMQ

    RabbitMQ <-->|video.transcript.ready| NLP
    NLP -->|video.analysis.complete| RabbitMQ
    RabbitMQ -->|video.analysis.complete| Orchestrator

    style Frontend fill:#42b883
    style Orchestrator fill:#9b42f5
    style VideoProc fill:#00add8
    style NLP fill:#3776ab
    style RabbitMQ fill:#ff6600
    style Postgres fill:#336791
    style Redis fill:#dc382d
```

## 📈 System Flow Overview

This sequence diagram illustrates the detailed, step-by-step flow of a request through the system.

```mermaid
sequenceDiagram
    participant User
    participant Frontend as Vue.js Frontend
    participant Orchestrator as Core Orchestrator (C#)
    participant RabbitMQ
    participant VideoProcessor as Video Processor (Go)
    participant YouTube
    participant WhisperAPI as OpenAI Whisper API
    participant NlpAnalyzer as NLP Analyzer (Python)
    participant Database as PostgreSQL

    User->>Frontend: Paste YouTube URL
    Frontend->>Orchestrator: Submit analysis request
    Orchestrator->>RabbitMQ: Publish `video.download.request`
    RabbitMQ-->>VideoProcessor: Consume `video.download.request`
    VideoProcessor->>YouTube: Download audio via yt-dlp
    YouTube-->>VideoProcessor: Audio stream/file
    VideoProcessor->>WhisperAPI: Request transcription
    WhisperAPI-->>VideoProcessor: Return transcript text
    VideoProcessor->>RabbitMQ: Publish `video.transcript.ready`
    RabbitMQ-->>NlpAnalyzer: Consume `video.transcript.ready`
    NlpAnalyzer->>NlpAnalyzer: Perform NLP + Sentiment Analysis
    NlpAnalyzer->>RabbitMQ: Publish `video.analysis.complete`
    RabbitMQ-->>Orchestrator: Consume `video.analysis.complete`
    Orchestrator->>Database: Store analysis results
    User->>Frontend: View analysis results
    Frontend->>Orchestrator: Fetch results
    Orchestrator->>Database: Query results
    Database-->>Orchestrator: Return results
    Orchestrator-->>Frontend: Return results to display
```

## 📨 RabbitMQ Event Model

| Event                     | Producer | Consumer | Description                                |
| ------------------------- | -------- | -------- | ------------------------------------------ |
| `video.transcript.requested`  | C#       | Go       | Request to fetch & transcribe video        |
| `video.transcript.completed`  | Go       | Python   | Transcript ready for NLP                   |
| `video.analysis.completed` | Python   | C#       | Final sentiment analysis ready for storage |

## 🔄 Data Flow Diagram

This diagram illustrates the data transformation pipeline:

```mermaid
flowchart LR
    A[📹 YouTube Video URL] --> B[🎵 Audio Extraction]
    B --> C[📝 Raw Transcript]
    C --> D[🔍 Feature Detection]
    D --> E[💭 Sentiment Analysis]
    E --> F[📊 Structured Results]

    subgraph "Video Processing Service"
        B
        C
    end

    subgraph "NLP Analyzer Service"
        D
        E
    end

    subgraph "Core Orchestrator"
        F
    end

    F --> G[(PostgreSQL)]
    F --> H[📱 Frontend Display]

    style A fill:#ff0000
    style C fill:#ffeb3b
    style F fill:#4caf50
    style G fill:#336791
    style H fill:#42b883
```

## 📊 Technology Stack Overview

| Layer                | Technologies                                          |
| -------------------- | ----------------------------------------------------- |
| **Frontend**         | Vue.js 3, Vite, Tailwind CSS v4, Pinia, Fetch API     |
| **Backend API**      | C#, ASP.NET Core, Entity Framework Core               |
| **Video Processing** | Go, Gin Framework, yt-dlp, FFmpeg, OpenAI Whisper API |
| **NLP & ML**         | Python                                                |
| **Message Queue**    | RabbitMQ 4.1.4 with Management Plugin                 |
| **Databases**        | PostgreSQL 18, Redis (planned)                        |
| **DevOps**           | Docker, Docker Compose                                |
| **External APIs**    | YouTube (via yt-dlp), OpenAI Whisper API              |

## 🤝 Contributing

Each service has its own README with detailed information:

- [Frontend Development](./core_frontend/README.md)
- [Backend Orchestrator](./core_orchestrator/README.md)
- [Video Processing](./video_processing_service/README.md)
- [NLP Analyzer](./NLP_analyzer_service/README.md)

## 📄 License

TBD
