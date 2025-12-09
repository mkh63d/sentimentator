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
                feat_conf = min(1.0, denom / (denom + 2.0))   # heuristic

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
    review = {'camera': ['This is a picture that I took on the S25.', 'And damn, that selfie camera looks so much better on Samsung.', 'But still, the photos come out looking more real.', 'They take incredible photos, they never lag, and somehow, they feel soulless.', 'You took a photo, it looked bad and you posted it anyway.', 'What shocked me even more was when I actually tested the camera.', "I was then looking at the back of the phone, and I thought, There's no way that tiny camera is still good.", 'And whoa, the selfie camera is so shaky.', 'And a picture of my dad, I love my dad.', "I didn't even know this, but even back then, you could just adjust the brightness of the torch, and on camera, the selfie looks like a Terminator.", 'And one last time when I tested the camera, I noticed that the phone does use HDR.', 'First of all, nowadays, the HDR effect is crazy.'], 'performance': ['They take incredible photos, they never lag, and somehow, they feel soulless.', 'And Google Maps was so slow that it would actually help you get more lost than you were before.', 'The scrolling was a little bit laggy.', 'No lag whatsoever.'], 'design': ['A beloved design, now with more love, whatever that means.', 'The design is crazy, the rose gold is gorgeous.'], 'speaker': ["The bottom has the fattest speaker grills I have ever seen, and in case there wasn't enough holes, there's also the headphone jack."], 'screen': ['Why do the colors look so good?', 'But then I made it to the home screen, and ah, the sweet memories, iOS 15.', "I didn't even know this, but even back then, you could just adjust the brightness of the torch, and on camera, the selfie looks like a Terminator."], 'charging': ['The night before, I made sure to charge my phone with the good old lightning cable, and... That sounds dusty.', "And I don't wanna spoil too much, but that's literally the first thing that I did after a full charge.", "Of course, at this point, iPhones didn't even have wireless charging yet."], 'battery': ['But I was pretty concerned to see only 9% of my battery left after just one mission.', 'Anyway, as good as the phone was, the battery died way too quickly before I could play any games or even listen to some music.', 'And so was the battery.']}
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

