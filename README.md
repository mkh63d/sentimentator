# 🎥 Sentimentator - Review Sentiment Analyzer

## 🧠 Overview

**Sentimentator** is a distributed, event-driven system that automatically analyzes **video reviews** of smartphones (from YouTube, ~~TikTok~~, and ~~Instagram~~) and produces a **synthetic sentiment summary** of opinions about device features (e.g., camera, battery, screen, performance).

## 🌍 High-Level Architecture

```mermaid
flowchart TD
    subgraph Frontend [🌐 Vue.js Frontend]
        A[User Interface]
    end

    subgraph CSharp [🏗️ C# Orchestrator]
        B1[Job Manager]
        B2[Database Handler]
        B3[Stuff]
        B4[More Stuff]
    end

    subgraph GoService [🦦 Go Downloader + Transcriber]
        C1[Video Downloader]
        C2[Speech-to-Text Engine]
    end

    subgraph PyService [🐍 Python NLP & Sentiment Analyzer]
        D1[NLP Processor]
        D2[Feature Sentiment Engine]
    end

    subgraph Infra [☁️ Infrastructure]
        E1[(RabbitMQ)]
        E2[(PostgreSQL)]
    end

    A -->|Video Link Request| B1
    B1 -->|video.download.request| E1
    E1 -->|Consume| C1
    C1 -->|Extract + Transcribe| C2
    C2 -->|video.transcript.ready| E1
    E1 -->|Consume| D1
    D1 -->|Analyze + Classify| D2
    D2 -->|video.analysis.complete| E1
    E1 -->|Consume| B2
    B2 -->|Store Results| E2
    B2 -->|Send Results| A
```

# System Flow Overview

```mermaid
sequenceDiagram
    participant User
    participant Vue as Vue.js Frontend
    participant CSharp as CSharp Orchestrator
    participant MQ as RabbitMQ
    participant Go as Go Downloader + Transcriber
    participant Py as Python NLP
    participant DB as PostgreSQL

    User->>Vue: Paste video URL
    Vue->>CSharp: Submit analysis request
    CSharp->>MQ: Publish "video.download.request"
    Go->>MQ: Consume "video.download.request"
    Go->>Go: Download video + extract audio + transcribe text
    Go->>MQ: Publish "video.transcript.ready"
    Py->>MQ: Consume "video.transcript.ready"
    Py->>Py: Perform NLP + Sentiment + Feature extraction
    Py->>MQ: Publish "video.analysis.complete"
    CSharp->>MQ: Consume "video.analysis.complete"
    CSharp->>DB: Store results
    Vue->>CSharp: Fetch results for visualization
```

## 📨 RabbitMQ Event Model

| Event                     | Producer | Consumer | Description                                |
| ------------------------- | -------- | -------- | ------------------------------------------ |
| `video.download.request`  | C#       | Go       | Request to fetch & transcribe video        |
| `video.transcript.ready`  | Go       | Python   | Transcript ready for NLP                   |
| `video.analysis.complete` | Python   | C#       | Final sentiment analysis ready for storage |
