# Action Plan: Week 4 – Planning and Research

## Tasks

1. **Precisely define the project scope (MVP - Minimum Viable Product).**
2. **Analyze available tools and APIs for downloading data from YouTube.**

<details>
<summary>2. Tools & APIs for Downloading Data from YouTube</summary>

---

#### 2. Tools & APIs for Downloading Data from YouTube

**1. YouTube Data API v3**

-   Use Cases:
    -   Search for videos
    -   Access video metadata (titles, descriptions, tags, etc.)
    -   Fetch captions ([Docs](https://developers.google.com/youtube/v3/guides/implementation/captions#captions-download))
-   Limitations:
    -   Quota limits
    -   Not all videos have captions available
    -   Cannot download video/audio content directly

**2. yt-dlp / go-ytdlp**

-   Use Cases:
    -   Download videos and audio from YouTube
-   Limitations:
    -   Subject to YouTube's terms of service
    -   Requires local processing and storage
-   Example Usage: - Download subtitles only (e.g., English):

    ```powershell
    .\yt-dlp.exe --sub-lang en --write-sub --skip-download https://www.youtube.com/watch?v=0e3GPea1Tyg
    ```

-   Download audio only:
    ```powershell
    .\yt-dlp.exe -x --audio-format mp3 https://www.youtube.com/watch?v=8HzW94z0MjU
    ```

---

</details>
<br/>

3. **Research and compare Speech-to-Text models**
    - Test and compare Whisper vs. Google API.

<details>
<summary>3. Speech-to-Text Models: Whisper vs. Google Cloud Speech-to-Text API</summary>

---

#### 3. Speech-to-Text Models Comparison

**_NOTE:_** Test TS scripts can be found in the `action_plans/week-4/resources`

**OpenAI Whisper API (whisper-1)**

**Pros:**

-   ✅ **Simpler implementation**: Single API call, no infrastructure setup required
-   ✅ **No additional storage needed**: Direct file upload (no GCS bucket required)
-   ✅ **Better punctuation**: Superior automatic punctuation and formatting
-   ✅ **Auto language detection**: Automatically detects language without configuration
-   ✅ **Fewer dependencies**: Only requires `openai` package

**Cons:**

-   ❌ **File size limit**: 25 MB maximum audio file size
-   ❌ **Limited customization**: Fewer configuration options
-   ❌ **Cost per minute**: $0.006/minute ($0.36/hour)

**Transcription Quality:**

```
It's hot today, isn't it? Yes, it is. I can't wait to get inside with the
air conditioning. Same here. Are you here for the conference? Yes, I'm hoping
to meet some new contacts. And yourself? Yes, same here. My name is Matt
Jenkins, by the way. Nice to meet you, Matt. My name is Olivia Christie.
So, what company are you with, Olivia? I'm working as a client relations
manager for Sun Systems. Ah, yes. You've just launched your new product,
right? Yes, so it's a very exciting time for us. Which company do you work
for, Matt? I work for Star Corporation as a sales manager. Oh, look. Here
comes the bus. Great. Let's enjoy the air conditioning.
```

---

**Google Cloud Speech-to-Text API**

**Pros:**

-   ✅ **No file size limit**: Can handle larger audio files via GCS if GCP storage is used
-   ✅ **Highly customizable**: Multiple configuration options (sample rate, encoding, etc.)
-   ✅ **Free tier**: 60 minutes/month free transcription

**Cons:**

-   ❌ **Slower processing**: Because for clips longer than 1 minute, requires GCS upload and GCP URL referencing
-   ❌ **Complex setup**: Requires GCS bucket creation and file upload
-   ❌ **More dependencies**: Requires `@google-cloud/speech` and `@google-cloud/storage`
-   ❌ **Authentication complexity**: Requires service account JSON key management, why not just API key?
-   ❌ **Worse punctuation**: Inconsistent punctuation (missing periods, commas)
-   ❌ **Manual language config**: Must specify language code (e.g., "en-US")

**Transcription Quality:**

```
It's hot today, isn't it?  Yes, it is. I can't wait to get inside with the
air conditioning.  Are you here for the conference? Yes, I'm hoping to meet
some new contacts and yourself. Yes, same here. My name is Matt Jenkins by
the way. Nice to meet you. Matt. My name is Olivia Christy.  So what company
are you with Olivia?  I'm working as a client relations manager for sun
systems. Yes, you've just launched your new product, right? Yes. So it's a
very exciting time for us.  Which company do you work for Matt? I work for
Star Corporation as a sales manager. Oh look. Here comes the bus. Let's
enjoy the air conditioning.
```

---

**Cost Analysis (Estimated for MVP)**

Assuming 100 hours of audio transcription per month:

| Service        | Pricing Model       | Monthly Cost (100 hrs) | Free Tier             |
| -------------- | ------------------- | ---------------------- | --------------------- |
| **Whisper**    | $0.006/minute       | $36.00                 | None                  |
| **Google STT** | $0.006-$0.024/min\* | $36.00-$144.00         | 60 min/month (~$0.36) |

\* [Google pricing doc](https://cloud.google.com/speech-to-text/pricing?hl=en)

\* [OpenAI Whisper Pricing Calculator](https://invertedstone.com/calculators/whisper-pricing)

---

**📋 Recommendation & Statement**

**For the Sentimentator MVP, OpenAI Whisper API is marked as winner**

1. **Development Velocity**:

    - Simpler implementation
    - No cloud storage bucket management required

2. **Performance**:

    - Higher quality transcription with better punctuation
    - Handles noise in audio sapless better
    - Automatic language detection eliminates configuration overhead

3. **Cost Efficiency for MVP**:

    - Same base cost as Google ($0.006/minute)
    - No hidden costs for storage or data transfer
    - Predictable pricing model
    - No free tier 😢

4. **Lower Operational Complexity**:

    - Single dependency (`openai` package)
    - Direct file upload without intermediate storage
    - Simpler credential management (API key vs. service account JSON)

---

</details>
<br/>

4. **Select and justify the technology stack**
    - Choose languages, frameworks, and databases.
5. **Create project documentation**
    - Vision document
    - Test plan
    - Division of tasks within the team
    - Preparation of the development environment
    - Report on the review of available options

## Goal

-   Deliver comprehensive project documentation covering the above points.
-   Ensure the team is ready for development with a clear plan and environment setup.

```

```

```

```
