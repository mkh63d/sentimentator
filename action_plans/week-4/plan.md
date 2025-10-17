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
