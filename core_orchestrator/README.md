# Sentimentator: C# Backend Orchestrator Service

## 🎯 Purpose

Central service responsible for:

-   Accepting frontend requests to analyze videos.
-   Publishing jobs to RabbitMQ (`video.download.request`).
-   Listening for analysis results (`video.analysis.complete`).
-   Storing and serving results from PostgreSQL.
-   ~~Notifying frontend (via server-sent events) when analysis is available.~~

---

### **Key Deliverables**

```mermaid
mindmap
  root((C# Backend))
    API Layer
      REST Endpoints
      OpenAPI/Swagger
      Input Validation
      Error Handling
    Authentication
      JWT Tokens???
      User Management
    Job Orchestration
      Workflow Management
      Status Tracking
      Error Recovery
      Progress Monitoring
    Data Layer
      Entity Framework
      Database Migrations
      Data Validation
    Integration
      Message Queue Client
      External API Calls
      Caching Strategy
      Background Services
```

## 🛠️ Technology Stack

| Component         | Technology          | Purpose                     |
| ----------------- | ------------------- | --------------------------- |
| **Framework**     | ???                 | Web API framework           |
| **Database**      | ???                 | Data persistence            |
| **Cache**         | Redis               | Performance optimization    |
| **Queue**         | RabbitMQ            | Inter-service communication |
| **Auth**          | JWT Bearer          | Authentication              |
| **Documentation** | Swagger/OpenAPI ??? | API documentation           |
| **Testing**       | xD                  | Unit/Integration tests      |
| **Monitoring**    | ???                 | Logging and monitoring      |

```mermaid
sequenceDiagram
    participant Client as Frontend Client
    participant API as C# API Gateway
    participant Queue as Message Queue
    participant DB as Database

    Client->>API: POST /api/videos/analyze
    Note over Client,API: {"url": "https://youtube.com/watch?v=xyz"}

    API->>DB: Create VideoAnalysisJob
    API->>Queue: Publish video.download message
    API->>Client: Return Job ID & Status

    loop Job Processing
        Queue-->>API: Status updates from Go/Python services
        API->>DB: Update job status
        API-->>Client: WebSocket notification
    end

    Client->>API: GET /api/results/{jobId}
    API->>DB: Retrieve analysis results
    API->>Client: Return complete sentiment analysis
```

```mermaid
graph LR
    subgraph "Outbound Messages"
        VDQ[video.download.queue<br/>→ Go Service]
        TAQ[text.analyze.queue<br/>→ Python Service]
    end

    subgraph "Inbound Messages"
        SUQ[status.update.queue<br/>← All Services]
        RRQ[results.ready.queue<br/>← Python Service]
    end

    subgraph "C# Message Handlers"
        JobOrch[Job Orchestrator]
        StatusHandler[Status Update Handler]
        ResultHandler[Result Processing Handler]
    end

    VDQ --> JobOrch
    TAQ --> JobOrch
    SUQ --> StatusHandler
    RRQ --> ResultHandler
```
