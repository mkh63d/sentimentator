import fs from "node:fs";
import { performance } from "perf_hooks";

import OpenAI from "openai";

import dotenv from "dotenv";
dotenv.config();

const AUDIO_FILE = process.env.AUDIO_FILE || "./speech_to_text_test_audio.mp3";
const OUTPUT_TEXT = "whisper_stt_transcription_test.txt";
const OUTPUT_META = "whisper_stt_metadata.json";

console.log("Starting Whisper API test...");
const openai = new OpenAI({ apiKey: process.env.OPENAI_API_KEY });

const start = performance.now();

console.log("Sending request to OpenAI Whisper API...");
const file = fs.createReadStream(AUDIO_FILE);

const transcription = await openai.audio.transcriptions.create({
    file: file,
    model: "whisper-1",
});

const duration = ((performance.now() - start) / 1000).toFixed(2);

fs.writeFileSync(OUTPUT_TEXT, transcription.text.trim(), "utf8");

const metadata = {
    model: "whisper-1",
    language: "auto-detected",
    processingTimeSeconds: duration,
    created: new Date().toISOString(),
    audioFile: AUDIO_FILE,
    outputFile: OUTPUT_TEXT,
    openaiResponse: {
        requestId: transcription._request_id ?? "not available",
        textLength: transcription.text?.length ?? "not available",
        usage: transcription.usage ?? "not available",
    },
};

fs.writeFileSync(OUTPUT_META, JSON.stringify(metadata, null, 2), "utf8");

console.log(`✅ Transcript saved to ${OUTPUT_TEXT}`);
console.log("Whisper API test completed.");
