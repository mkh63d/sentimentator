# language_preprocessor.py
import os
from langdetect import detect, DetectorFactory
from deep_translator import GoogleTranslator

# Ensure consistent language detection results
DetectorFactory.seed = 0


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

    if lang == "en":
        return text

    if lang == "unknown":
        return text

    try:
        translated = GoogleTranslator(source="auto", target="en").translate(text)
        return translated
    except Exception:
        return text


# ----------------------------------------------------------------------
# Standalone testing logic
# ----------------------------------------------------------------------
if __name__ == "__main__":
    # FIXED: go one dir up, then into resources/
    base_dir = os.path.dirname(__file__)               # .../src/
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

    # Detect language
    detected = detect_lang(text)
    print(f"[*] Detected Language: {detected}\n")

    # Process/translate
    processed = normalize_to_english(text)

    print("----- AFTER PROCESSING (ENGLISH) -----")
    print(processed)
    print("--------------------------------")
