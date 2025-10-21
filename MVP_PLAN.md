# MVP Plan for Sentimentator

## Overall MVP Goals

The goal of the MVP is to deliver a functional end-to-end application that allows users to submit a short YouTube video (which is about recension about mobile phone) URL and view a sentiment analysis report about mobile phone in it.

- **End-to-End Functionality:** A user can submit a short YouTube video URL via web platform and see the sentiment analysis results for mobile phone presented in it.
- **Core Features:** The system will extract key features (camera, battery, etc.) and determine the sentiment for each.
- **Basic UI:** A simple web interface to submit a URL and view the results.
- **Backend:** Backend that can handle the video processing and analysis workflow.

## Team-Specific Milestones

### Frontend (Vue.js) Team

- **Milestone 1: Basic UI for URL Submission:**
  - Create a simple form to accept a YouTube video URL.
  - Implement basic input validation, try to prevent submitting for example Reddit URLs.
- **Milestone 2: Results Display:**
  - Create a component to display the sentiment analysis results.
- **Milestone 3: API Integration:**
  - Integrate with the C# backend to allow users to:
    - login/register the account
    - submit allowed URLs
    - fetch and display results.
    - Handle loading and error states.

### Backend Orchestrator (C#) Team

- **Milestone 1: API Endpoints:**
  - Create REST endpoints for submitting a video URL and fetching results.
  - Implement OpenAPI/Swagger for API documentation.
- **Milestone 2: RabbitMQ Integration:**
  - Publish `video.download.request` messages when a new URL is submitted.
  - Consume `video.analysis.complete` messages and store the results.
- **Milestone 3: Database Integration:**
  - Set up a PostgreSQL database with Entity Framework.

### Video Processing (Go) Team

- **Milestone 1: Dev interface to input video URL:**
  - Create a simple web interface to input a YouTube video URL for testing purposes.
- **Milestone 2: RabbitMQ Integration:**
  - Consume `video.download.request` messages from RabbitMQ.
- **Milestone 3: Video Downloading and Audio Extraction:**
  - Use `yt-dlp` to download YouTube videos and extract audio from the downloaded videos.
- **Milestone 4: Transcription:**
  - Transcribe the audio to text.
  - Publish a `video.transcript.ready` event with the transcript.

### NLP Analyzer (Python) Team

- **Milestone 1: RabbitMQ Integration:**
  - Dev interface to input transcript text for testing purposes.
- **Milestone 2: RabbitMQ Integration:**
  - Consume `video.transcript.ready` messages from RabbitMQ.
- **Milestone 3: Sentiment Analysis:**
  - Use a pre-trained sentiment analysis model to analyze the transcript.
  - Identify key features (camera, battery, screen, etc.) using simple keyword matching.
- **Milestone 4: Publish Results:**
  - Structure the analysis results (overall sentiment, sentiment per feature).
  - Publish a `video.analysis.complete` event with the results.
