# feature_detection.py
import spacy
from typing import List, Dict

class FeatureDetector:
    """
    Detects smartphone-related features from text using spaCy.
    """

    DEFAULT_FEATURES = ["camera", "battery", "display", "speaker", "performance", "screen", "charging", "design"]

    def __init__(self, custom_features: List[str] = None):
        self.features = set(custom_features or self.DEFAULT_FEATURES)
        self.nlp = spacy.load("en_core_web_sm", disable=["ner", "parser"])
        if "sentencizer" not in self.nlp.pipe_names:
            self.nlp.add_pipe("sentencizer")

    def detect_features(self, text: str) -> Dict[str, List[str]]:
        """
        Returns dictionary of detected features with sentences mentioning them.
        Example:
        {
            "camera": ["The camera is amazing", "I love the camera quality"],
            "battery": ["The battery drains fast"]
        }
        """
        doc = self.nlp(text)
        sentences = [sent.text.strip() for sent in doc.sents]
        feature_mentions = {}

        for sent in sentences:
            lower_sent = sent.lower()
            for feat in self.features:
                if feat in lower_sent:
                    feature_mentions.setdefault(feat, []).append(sent)

        return feature_mentions


if __name__ == "__main__":
    text = """ The camera quality is absolutely amazing, the pictures are fantastic and stunning, but sometimes the autofocus fails.
            Battery life is good for light usage, though it drains quickly when gaming and under heavy load. 
            I love the bright and sharp display. 
            The speaker is loud and clear, but sometimes it distorts at high volume. 
            Performance is excellent for daily tasks, but lags with large apps.
            """
    fd = FeatureDetector()
    results = fd.detect_features(text)
    print(results)
