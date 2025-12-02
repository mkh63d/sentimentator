# feature_detection.py
import spacy
from typing import List, Dict

class FeatureDetector:
    """
    Detects smartphone-related features from text using spaCy.
    """

    DEFAULT_FEATURES = [
    # Original core features
    "camera",
    "battery",
    "display",
    "speaker",
    "performance",
    "screen",
    "charging",
    "design",

    # Hardware & Build
    "build_quality",
    "weight",
    "thickness",
    "water_resistance",
    "durability",

    # Software & Performance
    "os",
    "ui_ux",
    "app_compatibility",
    "gaming_performance",
    "lag",
    "updates",

    # Camera-specific
    "front_camera",
    "low_light_camera",
    "video_recording",
    "zoom",
    "stabilization",
    "color_accuracy_camera",

    # Battery & Charging
    "fast_charging",
    "wireless_charging",
    "battery_life_standby",
    "battery_health",

    # Display specifics
    "brightness",
    "refresh_rate",
    "color_accuracy_display",
    "screen_protection",
    "contrast",

    # Audio
    "microphone",
    "headphone_jack",
    "speaker_volume",
    "audio_clarity",

    # Connectivity
    "wifi",
    "bluetooth",
    "cellular",
    "gps",
    "nfc",

    # Practical & User Experience
    "storage",
    "security",
    "accessories",
    "price_value",
    "eco_friendliness"
]


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
    text = """So Samsung finally decided to upgrade a feature on the Galaxy S26 Ultra that's been stuck in 2017. Seven long years later it's finally happening. And that's not all, there are some pretty interesting updates on the phone too so let's get into it. What's up guys, Sal here so we already know Samsung is bumping the S26 Ultra's charging speed from 45W to 60W which is nice but what about wireless charging? What about MagSafe-style magnets or proper Cheeto support? Well, a new report from a reliable Korean outlet clears that up. Apparently, every phone in the Galaxy S26 lineup will get faster wireless charging but not all at the same level. The S26 Ultra is getting a big jump from 15W to 25W wireless charging. The regular S26 and S26 Plus will go to 20W compared to 15W before. So basically, wireless charging the S26 Ultra fully takes around 60 minutes compared to 100 minutes that 15W charging took. More importantly, all three phones will support Cheeto 2.2 which means they'll finally have built-in magnets inside the phone. That unlocks true Cheeto charging and lets you use proper magnetic accessories. Personally, I've been waiting for this for years. I love wireless charging and MagSafe-style accessories just make life easier. It's about time Samsung caught up. Now here's something a little unexpected. The front camera cutout on the S26 Ultra is getting bigger. Not only larger than the S25 Ultra's punch hole but the biggest Samsung has ever used on an Ultra phone. Ice Universe even shared a comparison image and yeah, it's very noticeable. So why the downgrade? Actually, there's a good reason behind it. The front camera's field of view is getting wider. Jumping from 80 degrees on the S25 Ultra to 85 degrees on the S26 Ultra. So you'll fit more in the frame, especially for wide selfies. But despite this wider field of view, Ice Universe says the camera specs themselves remain the same as last year's model. And just to clarify, this isn't the actual S26 Ultra camera sample, they are just visual references to show the difference in focal length. And it's not just the front camera getting tweaked, the main and periscope zoom lenses are also seeing focal length changes. But in the opposite direction, they are getting slightly more zoomed in. The main lens goes from 23mm to 24mm and the 5x zoom lens from 115mm to 120mm. You probably won't notice much difference on the main lens but the periscope one should be noticeable with better background compression and a more natural blur effect. Next up, we've got our first look at the Galaxy S26+, courtesy of Onleaks and surprise surprise, it looks exactly how we thought it would. Basically, a larger version of the base model. The display size remains 6.7 inches and the thickness sets at 7.3mm making it the only phone in the lineup that didn't slim down this year. Interestingly, it's not as rounded as the S25+. That's intentional, Samsung wants all three models to share the same corner design for a more unified look. But of course, there's some drama, Onleaks claims Samsung made some last-minute changes to the base S26 making it thicker than expected, 7.24mm instead of the rumored 6.9mm. Ice Universe disagrees saying that's false. He claims Onleaks only had early CAD data while he's seen the final official specs. According to him, the phone still measures 6.9mm which matches what we've been hearing for months. So classic leak vs leak situation but personally, I'm leaning toward Ice Universe on this one. When he's that confident, he usually ends up being right. Yeah, that just happened. Apple dropped a piece of cloth called the iPhone pocket for $230. And no, it's not a smart cloth, not a charging pouch, not even a case. Just a glorified sock for your phone for $230. And in typical Apple fashion, they are saying they made it with 3D knit technology which is basically a fancy way of saying they used a sewing machine and added an extra zero to the price. I mean, I thought the Apple polishing cloth was peak absurdity but this, this might just top it. Anyway, there are even more camera features coming to the Galaxy S26 Ultra plus a couple of really good updates you'll want to see. Click right here to watch the full breakdown and I'll see you there.
"""
    fd = FeatureDetector()
    results = fd.detect_features(text)
    print(results)
