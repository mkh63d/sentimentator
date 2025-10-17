import fs from "node:fs";
import OpenAI from "openai";
import { performance } from "perf_hooks";

console.log("Starting Whisper API test...");
const openai = new OpenAI({ apiKey: process.env.OPENAI_API_KEY });
const audioPath = "./speech_to_text_test_audio.mp3";

const start = performance.now();

const file = fs.createReadStream(audioPath);

console.log("Sending request to OpenAI Whisper API...");
const transcription = await openai.audio.transcriptions.create({
    file: file,
    model: "whisper-1",
});

const duration = ((performance.now() - start) / 1000).toFixed(2);

fs.writeFileSync(
    "whisper_transcription_test.txt",
    transcription.text.trim(),
    "utf8"
);

const metadata = {
    model: "whisper-1",
    audioPath,
    processingTimeSeconds: duration,
    created: new Date().toISOString(),
    openaiResponse: {
        request_id: transcription._request_id ?? null,
        textLength: transcription.text?.length ?? null,
    },
};

fs.writeFileSync(
    "whisper_metadata.json",
    JSON.stringify(metadata, null, 2),
    "utf8"
);

console.log("\n--- Whisper API Transcription Result ---");
console.log(`Transcription: ${transcription.text.trim()}`);
console.log(`Processing Time: ${duration} seconds}`);
console.log("Transcription saved to whisper_transcription.txt");
console.log("Metadata saved to whisper_metadata.json");
