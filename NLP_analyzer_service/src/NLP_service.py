# NLP_service.py
import json
import os
import sys
import threading
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../src")))

from consume import RabbitConsumer
from publish import RabbitPublisher
from feature_detection import FeatureDetector
from sentiment_analyzer import SentimentAnalyzer
from language_preprocessor import detect_lang, normalize_to_english

DEFAULT_LOCAL = "amqp://guest:guest@localhost:5672/"
RABBITMQ_URL = os.getenv("RABBITMQ_URL", DEFAULT_LOCAL)

class SentimentAnalysisService:
    """
    Main service that combines feature detection and sentiment analysis with publishing.
    """
    
    def __init__(self, rabbit_url=None):
        self.rabbit_url = rabbit_url or RABBITMQ_URL
        self._running = False
        self._service_thread = None
        
        self.consumer = RabbitConsumer(self.rabbit_url)
        self.features = FeatureDetector()
        self.sentiments = SentimentAnalyzer()
        self.publisher = RabbitPublisher(self.rabbit_url)
        
    def analyze_and_publish(self):
        """Complete pipeline: analyze text and publish results."""
        print("[*] Starting background consumer...")
        # Ensure consumer is started
        self.consumer.start_background()
        self._running = True

        while self._running:
            # wait_for_message with short timeout so stop() can interrupt loop
            if not self.consumer.wait_for_message(timeout=1.0):
                continue

            # Try to retrieve the full latest message (thread-safe)
            latest = None
            try:
                latest = self.consumer.get_latest_message()  # recommended method on consumer
            except Exception:
                # fallback to old accessors if get_latest_message() doesn't exist
                latest = None

            # If get_latest_message() not available, build from older helpers
            if not latest:
                analysis_id = self.consumer.get_analysis_id()
                transcript = self.consumer.get_text()
                # best-effort: original_url / video_id may not have helpers; set to None
                video_id = getattr(self.consumer, "get_video_id", lambda: None)()
                original_url = getattr(self.consumer, "get_original_url", lambda: None)()
                latest = {
                    "analysis_id": analysis_id,
                    "transcript": transcript,
                    "video_id": video_id,
                    "original_url": original_url,
                }

            # Support upstreams that nest payload under "data"
            if latest.get("data") and isinstance(latest["data"], dict):
                data = latest["data"]
                # Prefer top-level fields, fall back to data.*
                analysis_id = latest.get("analysis_id") or data.get("analysis_id")
                transcript = latest.get("transcript") or data.get("transcript") or latest.get("text")
                video_id = latest.get("video_id") or data.get("video_id")
                original_url = latest.get("original_url") or data.get("original_url")
            else:
                analysis_id = latest.get("analysis_id")
                transcript = latest.get("transcript") or latest.get("text")
                video_id = latest.get("video_id")
                original_url = latest.get("original_url")

            if not analysis_id or not transcript:
                print("[!] Received invalid or incomplete message, skipping...")
                print(f"    Full message: {json.dumps(latest, indent=2)}")
                print(f"    analysis_id: {analysis_id!r}")
                print(f"    transcript: {transcript!r}")
                print(f"    video_id: {video_id!r}")
                print(f"    original_url: {original_url!r}")

                # If 'data' exists, log it separately
                if "data" in latest and isinstance(latest["data"], dict):
                    print(f"    Nested 'data' content: {json.dumps(latest['data'], indent=2)}")
                continue

            print(f"[*] Analyzing text for video: {analysis_id}")
            detected_lang = detect_lang(transcript)
            print(f"[*] Detected Language: {detected_lang}")
            transcript = normalize_to_english(transcript)

            try:
                # Feature detection: returns dict[feature] -> list[str]
                feature_results = self.features.detect_features(transcript)
                print(f"[*] Detected features: {list(feature_results.keys())}")

                # Get per-feature sentiments from analyzer
                per_feature = self.sentiments.analyze_feature_sentiments(feature_results)

                # If analyze_feature_sentiments already returns SentimentAnalysis shape,
                # use it directly. Otherwise, transform per-feature output into the
                # required shape (overall, features, confidence).
                if isinstance(per_feature, dict) and {"overall", "features", "confidence"} <= set(per_feature.keys()):
                    sentiment_payload = per_feature
                else:
                    # Build numeric feature scores (in range -1..1) and confidence
                    numeric_features = {}
                    confidences = []
                    for feat, metrics in per_feature.items():
                        # Expected metrics from current analyzer:
                        # {"positive_score": x, "negative_score": y, "neutral_score": z, "overall": "positive"/"negative"/"neutral"}
                        pos = float(metrics.get("positive_score", 0.0))
                        neg = float(metrics.get("negative_score", 0.0))
                        neu = float(metrics.get("neutral_score", 0.0))
                        denom = pos + neg + neu
                        if denom <= 0:
                            feat_score = 0.0
                            feat_conf = 0.0
                        else:
                            feat_score = (pos - neg) / denom  # -1..1 (positive -> +, negative -> -)
                            # heuristic confidence: how many votes relative to a small smoothing constant
                            feat_conf = min(1.0, denom / (denom + 5.0))
                        numeric_features[feat] = feat_score
                        confidences.append(feat_conf)

                    # overall is mean of feature scores (range -1..1)
                    if numeric_features:
                        overall = sum(numeric_features.values()) / len(numeric_features)
                        confidence = sum(confidences) / len(confidences)
                    else:
                        overall = 0.0
                        confidence = 0.0

                    sentiment_payload = {
                        "overall": overall,
                        "features": numeric_features,
                        "confidence": confidence,
                    }

                # Publish: publisher should accept explicit fields (see note below)
                self.publisher.publish_analysis(
                    analysis_id=analysis_id,
                    original_url=original_url,
                    video_id=video_id,
                    transcript=transcript,
                    sentiment_results=sentiment_payload,
                )
                print(f"[*] Published analysis for video {analysis_id}")

            except Exception as e:
                print(f"[*] Error processing video {analysis_id}: {e}")
                
    def start(self):
        """Start the service in a background thread."""
        self._service_thread = threading.Thread(target=self.analyze_and_publish, daemon=True)
        self._service_thread.start()
        return self._service_thread
        
    def stop(self):
        """Stop the service gracefully."""
        print("[*] Stopping service...")
        self._running = False
        
    def close(self):
        """Close all connections gracefully."""
        print("[*] Closing all connections...")
        self.stop()
        
        # Wait a bit for the loop to break
        time.sleep(1)
        
        # Close components
        self.publisher.close()
        self.consumer.stop()
        
        # Wait for thread to finish
        if self._service_thread and self._service_thread.is_alive():
            print("[*] Waiting for service thread to finish...")
            self._service_thread.join(timeout=5)
            if self._service_thread.is_alive():
                print("[-] Service thread did not finish gracefully")
        
        print("[*] Service shutdown complete")

if __name__ == "__main__":
    print("[*] Running full NLP pipeline demo...\n")

    # Load mock transcript
    base_dir = os.path.dirname(__file__)
    project_root = os.path.abspath(os.path.join(base_dir, ".."))
    transcript_path = os.path.join(project_root, "resources", "mock_transcript.txt")

    if not os.path.exists(transcript_path):
        print(f"[!] Missing transcript file: {transcript_path}")
        sys.exit(1)

    with open(transcript_path, "r", encoding="utf-8") as f:
        transcript = f.read()

    print("----- ORIGINAL TRANSCRIPT -----")
    print(transcript)
    print("--------------------------------\n")

    # Fake message metadata
    analysis_id = "demo_001"
    video_id = "mock_video"
    original_url = "http://example.com/mock"

    # 1. Detect language
    detected = detect_lang(transcript)
    print(f"[*] Detected Language: {detected}")

    # 2. Translate to English if needed
    processed_transcript = normalize_to_english(transcript)
    print("\n----- NORMALIZED TO ENGLISH -----")
    print(processed_transcript)
    print("---------------------------------\n")

    # 3. Run feature detection
    features = FeatureDetector()
    feature_results = features.detect_features(processed_transcript)
    print(f"[*] Detected features: {list(feature_results.keys())}")

    # 4. Run sentiment analysis
    sentiments = SentimentAnalyzer()
    per_feature = sentiments.analyze_feature_sentiments(feature_results)

    # Normalize output shape
    if isinstance(per_feature, dict) and {"overall", "features", "confidence"} <= set(per_feature.keys()):
        sentiment_payload = per_feature
    else:
        numeric_features = {}
        confidences = []

        for feat, metrics in per_feature.items():
            pos = float(metrics.get("positive_score", 0.0))
            neg = float(metrics.get("negative_score", 0.0))
            neu = float(metrics.get("neutral_score", 0.0))
            denom = pos + neg + neu

            if denom <= 0:
                feat_score = 0.0
                feat_conf = 0.0
            else:
                feat_score = (pos - neg) / denom
                feat_conf = min(1.0, denom / (denom + 5.0))

            numeric_features[feat] = feat_score
            confidences.append(feat_conf)

        overall = sum(numeric_features.values()) / len(numeric_features) if numeric_features else 0.0
        confidence = sum(confidences) / len(confidences) if confidences else 0.0

        sentiment_payload = {
            "overall": overall,
            "features": numeric_features,
            "confidence": confidence,
        }

    print("\n----- FINAL SENTIMENT RESULTS -----")
    print(json.dumps(sentiment_payload, indent=2))
    print("-----------------------------------\n")

    print("[*] Demo pipeline completed successfully.")
