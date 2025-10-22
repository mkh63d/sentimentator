import fs from "node:fs";
import path from "node:path";
import { performance } from "perf_hooks";

import { SpeechClient } from "@google-cloud/speech";
import { Storage } from "@google-cloud/storage";

import dotenv from "dotenv";
dotenv.config();

const AUDIO_FILE = process.env.AUDIO_FILE || "./speech_to_text_test_audio.mp3";
const OUTPUT_TEXT = "google_stt_transcription_test.txt";
const OUTPUT_META = "google_stt_metadata.json";
const LANGUAGE_CODE = process.env.LANGUAGE_CODE || "en-US";
const BUCKET_NAME = process.env.GCS_BUCKET || "stt-temp-bucket";

const credentials = JSON.parse(fs.readFileSync("./gcp_api_key.json", "utf8"));

const speechClient = new SpeechClient({
    credentials,
});
const storage = new Storage({ credentials });

console.log("🎙️ Starting Google Speech-to-Text API test...");

async function uploadToGCS(
    filePath: string,
    bucketName: string
): Promise<string> {
    const fileName = path.basename(filePath);
    const bucket = storage.bucket(bucketName);

    const [exists] = await bucket.exists();
    if (!exists) {
        console.log(`🪣 Bucket "${bucketName}" does not exist. Creating...`);
        await bucket.create({ name: bucketName });
    }

    console.log(`📤 Uploading ${fileName} to gs://${bucketName}...`);
    await bucket.upload(filePath, { destination: fileName });
    return `gs://${bucketName}/${fileName}`;
}

async function run() {
    const start = performance.now();

    try {
        const stats = fs.statSync(AUDIO_FILE);

        let request: any;

        const config = {
            encoding: "MP3" as const,
            sampleRateHertz: 16000,
            languageCode: LANGUAGE_CODE,
            enableAutomaticPunctuation: true,
        };

        console.log("🕒 Using long-running recognition with GCS URI...");
        const gcsUri = await uploadToGCS(AUDIO_FILE, BUCKET_NAME);
        request = { config, audio: { uri: gcsUri } };

        const [operation] = await speechClient.longRunningRecognize(request);
        const [response] = await operation.promise();
        const transcript =
            response.results
                ?.map((r: any) => r.alternatives[0].transcript)
                .join(" ")
                .trim() || "Something went wrong, no transcription available.";

        const duration = ((performance.now() - start) / 1000).toFixed(2);
        fs.writeFileSync(OUTPUT_TEXT, transcript, "utf8");

        const metadata = {
            model: "Google Speech-to-Text API",
            language: LANGUAGE_CODE,
            processingTimeSeconds: duration,
            created: new Date().toISOString(),
            audioFile: AUDIO_FILE,
            outputFile: OUTPUT_TEXT,
            gcpResponse: {
                requestId: response.requestId || "not available",
                totalResults: response.results?.length || "not available",
                gcsUri: gcsUri,
            },
        };
        fs.writeFileSync(
            OUTPUT_META,
            JSON.stringify(metadata, null, 2),
            "utf8"
        );
        console.log(`✅ Transcript saved to ${OUTPUT_TEXT}`);
    } catch (err) {
        console.error("❌ Error during recognition:", err);
    }
}

run();
