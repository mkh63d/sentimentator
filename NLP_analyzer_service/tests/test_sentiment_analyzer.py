# tests/test_sentiment_analyzer.py
import os
import sys
import pytest

# Add src folder to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../src")))

from feature_detection import FeatureDetector
from sentiment_analyzer import SentimentAnalyzer

RESOURCE_DIR = os.path.join(os.path.dirname(__file__), "../resources")
REVIEW_FILE = os.path.join(RESOURCE_DIR, "review_text.txt")

# Expected sentiment results based on documented review
EXPECTED_SENTIMENTS = {
    "camera": "positive",       # 2 positive, 1 negative → we decide “positive”
    "battery": "neutral",      # 1 positive, 2 negative → we decide “negative”
    "display": "positive",
    "speaker": "negative",      # 1 positive, 1 negative → we decide “negative”
    "performance": "positive",  # 1 positive, 1 negative → we decide “positive”
}

@pytest.mark.unit
def test_sentiment_analysis_complex_review():
    # Load review text
    with open(REVIEW_FILE, "r", encoding="utf-8") as f:
        review_text = f.read()

    # Step 1: detect features
    fd = FeatureDetector()
    features = fd.detect_features(review_text)

    # Step 2: analyze sentiment
    sa = SentimentAnalyzer()
    sentiments = sa.analyze_feature_sentiments(features)

    # Basic sanity check
    assert isinstance(sentiments, dict)

    # Verify sentiment correctness
    for feature, expected in EXPECTED_SENTIMENTS.items():
        assert feature in sentiments, f"Feature '{feature}' missing from output"
        overall = sentiments[feature]["overall"]
        assert overall in ["positive", "negative", "neutral"], f"Unexpected overall sentiment '{overall}' for '{feature}'"
        assert overall == expected, f"Sentiment mismatch for '{feature}': expected '{expected}', got '{overall}'"