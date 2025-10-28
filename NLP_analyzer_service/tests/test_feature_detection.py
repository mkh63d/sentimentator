# tests/test_feature_detection.py

import os
import sys
import pytest
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../src")))
from feature_detection import FeatureDetector

RESOURCE_DIR = os.path.join(os.path.dirname(__file__), "../resources")
REVIEW_FILE = os.path.join(RESOURCE_DIR, "review_text.txt")

@pytest.mark.unit
def test_feature_detection_complex_review():
    # Load review text
    with open(REVIEW_FILE, "r", encoding="utf-8") as f:
        review_text = f.read()

    fd = FeatureDetector()
    detected = fd.detect_features(review_text)

    # Sanity checks
    assert isinstance(detected, dict)
    assert "camera" in detected
    assert "battery" in detected
    assert "display" in detected
    assert "speaker" in detected
    assert "performance" in detected

    # Check expected sentences for camera
    camera_sents = detected["camera"]
    assert any("amazing" in s for s in camera_sents)
    assert any("autofocus" in s for s in camera_sents)

    # Check battery
    battery_sents = detected["battery"]
    assert any("good" in s for s in battery_sents)
    assert any("drains quickly" in s for s in battery_sents)

    # Check display
    display_sents = detected["display"]
    assert any("bright" in s for s in display_sents)

    # Check speaker
    speaker_sents = detected["speaker"]
    assert any("loud" in s for s in speaker_sents)
    assert any("distorts" in s for s in speaker_sents)

    # Check performance
    perf_sents = detected["performance"]
    assert any("excellent" in s for s in perf_sents)
    assert any("lags" in s for s in perf_sents)
