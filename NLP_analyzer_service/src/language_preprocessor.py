# language_preprocessor.py
import os
import re
from langdetect import detect, DetectorFactory
from deep_translator import GoogleTranslator

# Ensure consistent language detection results
DetectorFactory.seed = 0

def chunk_text(text: str, max_len: int = 4500):
    """
    Splits text into chunks < max_len while respecting sentence boundaries.
    """
    sentences = re.split(r'(?<=[.!?])\s+', text)
    chunks = []
    current = ""

    for s in sentences:
        # If sentence itself is longer than max_len, hard-split it
        if len(s) > max_len:
            # break the long sentence into strict chunks
            for i in range(0, len(s), max_len):
                part = s[i:i + max_len]
                if current:
                    chunks.append(current)
                    current = ""
                chunks.append(part)
            continue

        # append sentence normally
        if len(current) + len(s) + 1 > max_len:
            chunks.append(current)
            current = s
        else:
            current = s if not current else current + " " + s

    if current:
        chunks.append(current)

    return chunks


def translate_long_text(text: str, target: str = "en") -> str:
    """
    Translates long text by splitting into chunks.
    """
    chunks = chunk_text(text)
    translated_chunks = []

    for ch in chunks:
        try:
            translated = GoogleTranslator(source="auto", target=target).translate(ch)
        except Exception as e:
            print("[Translation Error]", e)
            translated = ch  # fallback to original chunk
        translated_chunks.append(translated)

    return " ".join(translated_chunks)

def detect_lang(text: str) -> str:
    """
    Returns the detected language code (e.g., 'en', 'fr', 'de').
    Falls back to 'unknown' on error.
    """
    try:
        return detect(text)
    except Exception:
        return "unknown"

def normalize_to_english(text: str) -> str:
    """
    Detects language; if not English, translates to English.
    Returns the English text (original or translated).
    """
    lang = detect_lang(text)

    if lang in ["en", "unknown"]:
        return text

    try:
        # For long texts, switch to chunking
        if len(text) > 4500:
            return translate_long_text(text, target="en")

        # For short ones, direct call
        return GoogleTranslator(source="auto", target="en").translate(text)

    except Exception as e:
        print("[!] Translation failed:", e)
        return text

if __name__ == "__main__":
    base_dir = os.path.dirname(__file__)
    project_root = os.path.abspath(os.path.join(base_dir, ".."))
    resources_path = os.path.join(project_root, "resources", "mock_transcript.txt")

    if not os.path.exists(resources_path):
        print(f"[!] File not found: {resources_path}")
        exit(1)

    print("[*] Loading transcript from mock_transcript.txt...\n")

    with open(resources_path, "r", encoding="utf-8") as f:
        text = f.read()

    print("----- ORIGINAL TRANSCRIPT -----")
    print(text)
    print("--------------------------------\n")

    detected = detect_lang(text)
    print(f"[*] Detected Language: {detected}\n")

    processed = normalize_to_english(text)

    print("----- AFTER PROCESSING (ENGLISH) -----")
    print(processed)
    print("--------------------------------")
