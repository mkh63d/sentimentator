import re
from transformers import pipeline

class SentimentAnalyzer:
    def __init__(self):
        # Model fine-tuned on product reviews (star rating output)
        self.analyzer = pipeline(
            "sentiment-analysis",
            model="nlptown/bert-base-multilingual-uncased-sentiment"
        )

    @staticmethod
    def split_into_clauses(text: str):
        raw_clauses = re.split(r',\s*(?:but|though|although|however)\s*', text, flags=re.IGNORECASE)
        clauses = [c.strip() for c in raw_clauses if c.strip()]
        return clauses

    @staticmethod
    def label_to_sentiment(label: str):
        """Convert '1 star' -> negative, '2 star' -> negative, '3 star' -> neutral, '4-5 star' -> positive"""
        star = int(label[0])
        if star <= 2:
            return "negative"
        elif star == 3:
            return "neutral"
        else:
            return "positive"

    def analyze_feature_sentiments(self, feature_sentences: dict):
        feature_sentiments = {}

        for feature, sentences in feature_sentences.items():
            clauses = []
            for sent in sentences:
                clauses.extend(self.split_into_clauses(sent))

            results = self.analyzer(clauses)

            pos_score, neg_score, neutral_score = 0.0, 0.0, 0.0
            for res in results:
                sentiment = self.label_to_sentiment(res["label"])
                score = res["score"]
                if sentiment == "positive":
                    pos_score += score
                elif sentiment == "negative":
                    neg_score += score
                else:
                    neutral_score += score

            # Decide overall sentiment
            scores = {"positive": pos_score, "negative": neg_score, "neutral": neutral_score}
            if abs(pos_score - neg_score) < 0.05:
                overall = "neutral"
            else:
                overall = max(scores, key=scores.get)

            feature_sentiments[feature] = {
                "positive_score": pos_score,
                "negative_score": neg_score,
                "neutral_score": neutral_score,
                "overall": overall,
            }

        return feature_sentiments


# ---------------- Main ----------------
if __name__ == "__main__":
    review = {
        "camera": ["The camera quality is absolutely amazing, the pictures are fantastic and stunning, but sometimes the autofocus fails."],
        "battery": ["Battery life is good for light usage, though it drains quickly when gaming and under heavy load."],
        "display": ["I love the bright and sharp display. "],
        "speaker": ["The speaker quality is bad and additionally sometimes it distorts at high volume. "],
        "performance": ["Performance is excellent for daily tasks, but lags with large apps."]
    }

    sa = SentimentAnalyzer()
    sentiments = sa.analyze_feature_sentiments(review)

    print("\nFeature Sentiments:\n")
    for feature, scores in sentiments.items():
        print(f"{feature}: {scores}")
