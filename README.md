# 🎥 Sentimentator - Review Sentiment Analyzer

## 🧠 Overview

**Sentimentator** is a distributed, event-driven system that automatically analyzes **video reviews** of smartphones from **YouTube** and produces a **synthetic sentiment summary** of opinions about device features (e.g., camera, battery, screen, performance).

## 🌍 High-Level Architecture

The system is composed of several microservices that communicate asynchronously via a message queue.

```mermaid
graph TB
    subgraph L1["Client Layer"]
        User[👤 User]
        Browser[🌐 Web Browser]
        User --> Browser
    end

    subgraph L2["Frontend Layer"]
        direction TB
            Frontend[Vue.js Frontend
            • URL Submission
            • Results Display]
    end

    subgraph L3["API Gateway Layer"]
        direction TB
        Orchestrator[C# Core Orchestrator
        • REST API
        • Job Management
        • Auth & Validation]
    end

    subgraph L4["Message Queue Layer"]
        direction TB
        RabbitMQ[🐰 RabbitMQ<br/>Port: 5672<br/>Management: 15672<br/>• Event Routing<br/>• Async Communication]
    end

    subgraph L5["Processing Services"]
        VideoProc[Go Video Processor<br/>• Content Audio Extraction<br/>• Transcript Generation]
        NLP[Python NLP Analyzer<br/>• Feature Detection<br/>• Sentiment Analysis]
    end

    subgraph L6["External APIs"]
        YouTube[📺 YouTube]
        Instagram[📺 Instagram]
        Whisper[🎙️ OpenAI Whisper API<br/>Speech-to-Text]
    end

    subgraph L7["Data Layer"]
        OrchestratorDB[(🗄️ PostgreSQL<br/>orchestrator DB<br/>──────────<br/>• Analysis<br/>• User & Identity)]
        VideoProcDB[(🗄️ PostgreSQL<br/>video_processing DB<br/>──────────<br/>• Content Metadata<br/>• Transcripts)]
    end

    Browser --> Frontend
    Frontend <-->|REST API| Orchestrator

    Orchestrator <-->|EF Core| OrchestratorDB
    Orchestrator -->|Publish: video.transcript.requested| RabbitMQ
    RabbitMQ -->|Consume: video.analysis.completed| Orchestrator

    RabbitMQ -->|Consume: video.transcript.requested| VideoProc
    VideoProc <-->|GORM| VideoProcDB
    VideoProc -->|Scrape| YouTube
    VideoProc -->|Scrape| Instagram
    VideoProc -->|Transcribe| Whisper
    VideoProc -->|Publish: video.transcript.completed| RabbitMQ

    RabbitMQ -->|Consume: video.transcript.completed| NLP
    NLP -->|Publish: video.analysis.completed| RabbitMQ

    style Frontend fill:#42b883
    style Orchestrator fill:#9b42f5
    style VideoProc fill:#00add8
    style NLP fill:#3776ab
    style RabbitMQ fill:#ff6600
    style OrchestratorDB fill:#336791
    style VideoProcDB fill:#336791
```

## 📈 System Flow Overview

This sequence diagram illustrates the detailed, step-by-step flow of a request through the system.

```mermaid
sequenceDiagram
    participant User
    participant Frontend as Vue.js Frontend
    participant Orchestrator as Core Orchestrator (C#)
    participant Exchange as RabbitMQ Exchange<br/>(sentimentator.topic)
    participant Q1 as Queue: video.processing
    participant VideoProcessor as Video Processor (Go)
    participant YouTube
    participant WhisperAPI as OpenAI Whisper API
    participant Q2 as Queue: video.transcript
    participant NlpAnalyzer as NLP Analyzer (Python)
    participant Q3 as Queue: video.analysis
    participant Database as PostgreSQL

    User->>Frontend: Paste YouTube URL
    Frontend->>Orchestrator: POST /analyses/{url}
    Orchestrator->>Database: Create Analysis record
    Orchestrator->>Exchange: Publish video.transcript.requested
    Exchange->>Q1: Route to video.processing
    Q1->>VideoProcessor: Consume message

    VideoProcessor->>YouTube: Fetch metadata + download audio via yt-dlp
    YouTube-->>VideoProcessor: Metadata + audio stream/file
    VideoProcessor->>Exchange: Publish video.metadata.completed
    Exchange->>Q4: Route to video.metadata
    Q4->>Orchestrator: Consume message
    Orchestrator->>Database: Update analysis with video metadata
    VideoProcessor->>WhisperAPI: Request transcription
    WhisperAPI-->>VideoProcessor: Return transcript text
    VideoProcessor->>Database: Store video & transcript
    VideoProcessor->>Exchange: Publish video.transcript.completed
    Exchange->>Q2: Route to video.transcript

    Q2->>NlpAnalyzer: Consume message
    NlpAnalyzer->>Database: Fetch transcript
    NlpAnalyzer->>NlpAnalyzer: Perform NLP + Sentiment Analysis
    NlpAnalyzer->>Exchange: Publish video.analysis.completed
    Exchange->>Q3: Route to video.analysis

    Q3->>Orchestrator: Consume message
    Orchestrator->>Database: Update analysis with results

    User->>Frontend: Poll for results
    Frontend->>Orchestrator: GET /analyses/{id}
    Orchestrator->>Database: Query results
    Database-->>Orchestrator: Return results
    Orchestrator-->>Frontend: Return analysis
    Frontend-->>User: Display sentiment results
```

## 📨 RabbitMQ Event Architecture

This diagram shows the detailed RabbitMQ topology with exchanges, queues, routing keys, and event flows:

```mermaid
graph TB
    subgraph "Publishers"
        Orchestrator[C# Core Orchestrator]
        VideoProc[Go Video Processor]
        NLP[Python NLP Analyzer]
    end

    subgraph "RabbitMQ Broker"
        Exchange[📮 Topic Exchange<br/>`sentimentator.topic`]

        subgraph "Queues"
            Q1[📦 video.processing<br/>Bound: video.transcript.requested]
            Q2[📦 video.transcript<br/>Bound: video.transcript.completed]
            Q3[📦 video.analysis<br/>Bound: video.analysis.completed]
            Q4[📦 video.metadata<br/>Bound: video.metadata.completed]
        end
    end

    subgraph "Consumers"
        VideoProcConsumer[Go Video Processor]
        NLPConsumer[Python NLP Analyzer]
        OrchestratorConsumer[C# Core Orchestrator]
    end

    Orchestrator -->|"Publish<br/>video.transcript.requested"| Exchange
    Exchange -->|"Route by key"| Q1
    Q1 --> VideoProcConsumer

    VideoProc -->|"Publish<br/>video.transcript.completed"| Exchange
    Exchange -->|"Route by key"| Q2
    Q2 --> NLPConsumer

    NLP -->|"Publish<br/>video.analysis.completed"| Exchange
    Exchange -->|"Route by key"| Q3
    Q3 --> OrchestratorConsumer

    VideoProc -->|"Publish<br/>video.metadata.completed"| Exchange
    Exchange -->|"Route by key"| Q4
    Q4 --> OrchestratorConsumer

    style Exchange fill:#ff6600,color:#fff
    style Q1 fill:#00add8,color:#fff
    style Q2 fill:#3776ab,color:#fff
    style Q3 fill:#9b42f5,color:#fff
    style Orchestrator fill:#9b42f5,color:#fff
    style VideoProc fill:#00add8,color:#fff
    style VideoProcConsumer fill:#00add8,color:#fff
    style NLP fill:#3776ab,color:#fff
    style NLPConsumer fill:#3776ab,color:#fff
    style OrchestratorConsumer fill:#9b42f5,color:#fff
```

### Event Details

| Event Name               | Routing Key                  | Publisher           | Queue              | Consumer            | Payload                                                                        |
| ------------------------ | ---------------------------- | ------------------- | ------------------ | ------------------- | ------------------------------------------------------------------------------ |
| **Transcript Requested** | `video.transcript.requested` | C# Orchestrator     | `video.processing` | Go Video Processor  | `{analysis_id, url, timestamp}`                                                |
| **Transcript Completed** | `video.transcript.completed` | Go Video Processor  | `video.transcript` | Python NLP Analyzer | `{analysis_id, video_id, transcript_id, text, language, duration, created_at}` |
| **Analysis Completed**   | `video.analysis.completed`   | Python NLP Analyzer | `video.analysis`   | C# Orchestrator     | `???`                                                                          |
| **Metadata Completed**   | `video.metadata.completed`   | Go Video Processor  | `video.metadata`   | C# Orchestrator     | `{analysis_id, video_id, title, thumbnail_url, duration_seconds?, ???}`        |

## 🔄 Data Flow Diagram

This diagram illustrates the data transformation pipeline:

```mermaid
flowchart LR
    A[📹 YouTube/Instagram Video URL] --> B[🎵 Audio Extraction]
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
| **Frontend**         | Vue.js 3, Vite, Tailwind CSS v4, Pinia, OpenAPI/fetch |
| **Backend API**      | C#, ASP.NET Core, Entity Framework Core               |
| **Video Processing** | Go, Gin, yt-dlp, FFmpeg, OpenAI Whisper API           |
| **NLP & ML**         | Python ???                                            |
| **Message Queue**    | RabbitMQ 4.1.4 with Management Plugin                 |
| **Databases**        | PostgreSQL 18                                         |
| **DevOps**           | Docker, Docker Compose                                |
| **External APIs**    | YouTube & Instagram (via yt-dlp), OpenAI Whisper API  |

## 🤝 Contributing

Each service has its own README with detailed information:

-   [Frontend Development](./core_frontend/README.md)
-   [Backend Orchestrator](./core_orchestrator/README.md)
-   [Video Processing](./video_processing_service/README.md)
-   [NLP Analyzer](./NLP_analyzer_service/README.md)

## 📄 License

TBD
