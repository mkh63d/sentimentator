# Sentimentator: C# Backend Orchestrator Service

## 🎯 Purpose

Central service responsible for:

-   Provide REST endpoints .
-   Publish `video.download.request` messages for the downloader service.
-   Consume `video.analysis.complete` messages and persist final structured results.
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
