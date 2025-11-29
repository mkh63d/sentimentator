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
        feature_confidences = []

        for feature, sentences in feature_sentences.items():
            clauses = []
            for sent in sentences:
                clauses.extend(self.split_into_clauses(sent))

            if not clauses:
                feature_sentiments[feature] = 0.0
                feature_confidences.append(0.0)
                continue

            results = self.analyzer(clauses)

            pos_score = neg_score = neu_score = 0.0
            for res in results:
                sentiment = self.label_to_sentiment(res["label"])
                score = float(res.get("score", 0.0))
                if sentiment == "positive":
                    pos_score += score
                elif sentiment == "negative":
                    neg_score += score
                else:
                    neu_score += score

            denom = pos_score + neg_score + neu_score
            if denom <= 0:
                feat_score = 0.0
                feat_conf = 0.0
            else:
                feat_score = (pos_score - neg_score) / denom  # -1..1
                feat_conf = min(1.0, denom / (denom + 5.0))   # heuristic

            feature_sentiments[feature] = feat_score
            feature_confidences.append(feat_conf)

        # overall: mean of feature scores
        if feature_sentiments:
            overall = sum(feature_sentiments.values()) / len(feature_sentiments)
            confidence = sum(feature_confidences) / len(feature_confidences)
        else:
            overall = 0.0
            confidence = 0.0

        return {
            "overall": overall,
            "features": feature_sentiments,
            "confidence": confidence,
        }



# ---------------- Main ----------------
if __name__ == "__main__":
    review = {
        'updates': [
            "And that's not all, there are some pretty interesting updates on the phone too so let's get into it.",
            "Anyway, there are even more camera features coming to the Galaxy S26 Ultra plus a couple of really good updates you'll want to see."
        ],
        'charging': [
            "What's up guys, Sal here so we already know Samsung is bumping the S26 Ultra's charging speed from 45W to 60W which is nice but what about wireless charging?",
            'Apparently, every phone in the Galaxy S26 lineup will get faster wireless charging but not all at the same level.',
            'The S26 Ultra is getting a big jump from 15W to 25W wireless charging.',
            'So basically, wireless charging the S26 Ultra fully takes around 60 minutes compared to 100 minutes that 15W charging took.',
            'That unlocks true Cheeto charging and lets you use proper magnetic accessories.',
            'I love wireless charging and MagSafe-style accessories just make life easier.',
            "And no, it's not  a smart cloth, not a charging pouch, not even a case."
        ],
        'accessories': [
            'That unlocks true Cheeto charging and lets you use proper magnetic accessories.',
            'I love wireless charging and MagSafe-style accessories just make life easier.'
        ],
        'camera': [
            'The front camera cutout on the S26 Ultra is getting bigger.',
            "The front camera's field of view is getting wider.",
            "But despite this wider field of view, Ice Universe says the camera specs themselves remain the same as last year's model.",
            "And just to clarify, this isn't the actual S26 Ultra camera sample, they are just visual references to show the difference in  focal length.",
            "And it's not just the front camera getting tweaked, the main and periscope zoom lenses are also seeing focal length changes.",
            "Anyway, there are even more camera features coming to the Galaxy S26 Ultra plus a couple of really good updates you'll want to see."
        ],
        'zoom': [
            "And it's not just the front camera getting tweaked, the main and periscope zoom lenses are also seeing focal length changes.",
            'But in the opposite direction, they are getting slightly more zoomed in.',
            'The main lens goes from 23mm to 24mm and the 5x zoom lens from 115mm to 120mm.'
        ],
        'os': ['But in the opposite direction, they are getting slightly more zoomed in.'],
        'thickness': [
            "The display size remains 6.7 inches and the thickness sets at 7.3mm making it the only phone in the lineup that didn't slim down this year."
        ],
        'display': [
            "The display size remains 6.7 inches and the thickness sets at 7.3mm making it the only phone in  the lineup that didn't slim down this year."
        ],
        'design': [
            "That's intentional, Samsung wants all three models to share the same corner design for a more unified look."
        ]
    }

    sa = SentimentAnalyzer()
    sentiments = sa.analyze_feature_sentiments(review)

    # If analyzer returned the SentimentAnalysis shape, use its 'features' map.
    if isinstance(sentiments, dict) and "features" in sentiments and isinstance(sentiments["features"], dict):
        features_map = sentiments["features"]
        overall = sentiments.get("overall")
        confidence = sentiments.get("confidence")
        print("\nOverall:", round(overall, 3) if isinstance(overall, (int, float)) else overall)
        print("Confidence:", round(confidence, 3) if isinstance(confidence, (int, float)) else confidence)
    else:
        # fallback to old shape (per-feature dict)
        features_map = sentiments

    print("\nFeature Sentiments:\n")
    for feature, score in features_map.items():
        # each score is expected to be a float in -1..1 (or any numeric)
        if isinstance(score, (int, float)):
            print(f"{feature}: {round(score, 3)}")
        else:
            # if older analyzer returned dict per feature, print that nicely
            try:
                clipped_scores = {
                    k: round(v, 3) if isinstance(v, (int, float)) else v
                    for k, v in score.items()
                }
                print(f"{feature}: {clipped_scores}")
            except Exception:
                print(f"{feature}: {score}")

