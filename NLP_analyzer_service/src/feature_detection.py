# feature_detection.py
import re
from typing import List, Dict
import spacy


class FeatureDetector:
    """
    Detects smartphone-related features from text using spaCy.
    """

    DEFAULT_FEATURES = [
        # Core features
        "camera", "battery", "speaker", "performance", "screen", "charging", "design",
        # Hardware & Build
        "build_quality", "weight", "thickness", "water_resistance", "durability",
        # Software & Performance
        "app_compatibility", "updates",
        # Audio
        "microphone", "headphone_jack", "speaker_volume", "audio_clarity",
        # Connectivity
        "wifi", "bluetooth", "cellular", "gps", "nfc",
        # Practical & User Experience
        "storage", "security", "accessories", "price_value", "eco_friendliness"
    ]

    DEFAULT_KEYWORD_MAPPING = {
        "camera": ["photo", "photos", "picture", "selfie", "photograph",
                   "stabilization", "video_recording", "color_accuracy", "zoom"],
        "battery": ["charge", "charging", "power", "battery life", "battery_health"],
        "screen": ["display", "screen", "resolution", "brightness",
                   "screen_protection", "contrast", "refresh_rate"],
        "speaker": ["audio", "speaker", "sound", "speaker_volume", "audio_clarity"],
        "performance": ["speed", "lag", "slow", "fast", "lag", "laggy", "gaming"],
        "storage": ["space"]
    }

    def __init__(self, custom_features: List[str] = None, keyword_mapping: Dict[str, List[str]] = None):
        self.features = set(custom_features or self.DEFAULT_FEATURES)

        # Load spaCy
        self.nlp = spacy.load("en_core_web_sm", disable=["ner", "parser"])
        if "sentencizer" not in self.nlp.pipe_names:
            self.nlp.add_pipe("sentencizer")

        # Merge default + custom keyword mappings
        combined_mapping = self.DEFAULT_KEYWORD_MAPPING.copy()
        if keyword_mapping:
            combined_mapping.update(keyword_mapping)

        # Compile keyword regex patterns
        self.keyword_patterns = {}
        for feat, keywords in combined_mapping.items():
            self.keyword_patterns[feat] = [
                re.compile(rf"\b{re.escape(kw)}\b", flags=re.IGNORECASE) for kw in keywords
            ]

        # Compile feature-name regex patterns
        self.feature_patterns = {
            feat: re.compile(rf"\b{re.escape(feat)}\b", flags=re.IGNORECASE)
            for feat in self.features
        }

    def detect_features(self, text: str) -> Dict[str, List[str]]:
        """
        Returns a dict of detected features with sentences mentioning them.
        Deduplicates sentences so they appear only once per feature.
        """
        doc = self.nlp(text)
        sentences = [sent.text.strip() for sent in doc.sents]

        # Use sets to avoid duplicates
        feature_mentions: Dict[str, set] = {}

        for sent in sentences:
            # Direct feature matches
            for feat, pattern in self.feature_patterns.items():
                if pattern.search(sent):
                    feature_mentions.setdefault(feat, set()).add(sent)

            # Keyword matches
            for feat, patterns in self.keyword_patterns.items():
                for pattern in patterns:
                    if pattern.search(sent):
                        feature_mentions.setdefault(feat, set()).add(sent)

        # Convert sets → lists for output
        return {feat: list(sents) for feat, sents in feature_mentions.items()}


if __name__ == "__main__":
    text = """Modern phones are perfect, too perfect. They take incredible photos, they never lag, and somehow, they feel soulless. So I went backwards, all the way back to the first iPhone SE. It's tiny, it's outdated, and honestly, it might be the last iPhone that actually felt human. But nostalgia lies, it tricks you into thinking that old equals good. So I spent a few days with this thing to find out if it's still good, or I just miss being happy. So I bought this phone used, off of CEX, and the description says, the most powerful 4-inch phone ever. I could say the same thing about myself, without the phone part. A beloved design, now with more love, whatever that means. But I'm dying to find out what the 4 inches actually feels like in my hand. As you can see, the box is massive, and it's completely beaten up. It's cooked, just like this table. I don't know why it's got so many stickers, but if we open it, wow. The boxes today are like, half the size of this, that's weird. And then of course we get a third party cable, and a plug. Apparently it's a luxury to have a charger these days. The design is crazy, the rose gold is gorgeous. The volume buttons are just two circles with a plus and a minus, how cool is that? The bottom has the fattest speaker grills I have ever seen, and in case there wasn't enough holes, there's also the headphone jack. I can't imagine a more beautiful thing. I can't imagine a more beautiful thing. I love it, it's beautiful. So I sat down in my jungle chair, yes, that's what I call it, to finally see what this piece of nostalgia is really like. And immediately, I did feel a sense of nostalgia. Um, my fingerprint scanner just decided to give up. But then I made it to the home screen, and ah, the sweet memories, iOS 15. How the hell did this phone get 6 years of software support? I have no idea. And speaking of the home button, this was the last phone to have a physical home button that you can click. Unlike the iPhone 7 that used vibrations to make it feel like you're clicking it. And why does that matter? The old notification panel looked pretty cool, until they changed it to the liquid ass, or I mean, the liquid glass update. And they made it look all weird and stretchy. I didn't even know this, but even back then, you could just adjust the brightness of the torch, and on camera, the selfie looks like a Terminator. And well, since the light is on, let me just show you my room really quick. I have a shopping cart full of boxes, the usual. A t-shirt hanging from the ceiling, at least it's not me. And a picture of my dad, I love my dad. And lastly, an unfinished puzzle under my bed. You're welcome, for the useless information. Alright, here's a fun fact. If you go to add an alarm, and scroll all the way down, you can actually get to the end, so it's not an infinite loop. And apparently, you can still do that. As I kept scrolling, I realised that this phone is built different. It's built from aluminium, don't tell me how to pronounce it. Which is the same as the iPhone 17 Pro. Which apparently gets scratches if you just look at it, which got downgraded from titanium. But is that a downgrade, or is that just going back to the old roots? Well, the real answer is... SE stands for Special Edition. As they took the insides of the iPhone 6S, and put it in the body of the iPhone 5S. Because apparently, a lot of people still preferred 4 inches. And apparently, so do I. Because when I tried it for a day, I was genuinely shocked at how good it felt. The night before, I made sure to charge my phone with the good old lightning cable, and... That sounds dusty. So today, I have an actual mission for this phone. I have a bag of random cables and chargers that I wanna recycle. Now for this mission, I will need Apple Pay to travel. I will need Google Maps, because I don't know where I'm going. And I'm not gonna use Apple Maps. I'm gonna need some entertainment for the journey. The most productive one being Instagram Reels. And maybe some games and Spotify. And most importantly, I'm gonna need this phone to last me more than 5 minutes. Can it do that? Well, let's find out. One of these phones has to go to the recycling center. An old HTC phone, or a Nokia 6300. And... Sorry HTC, but your time has come. I made sure to add my card to Apple Pay to make sure you could still even do that. And I don't wanna spoil too much, but that's literally the first thing that I did after a full charge. Google Maps installed perfectly fine. It didn't even care that I was using a 9-year-old phone. I set my destination and it was time to begin my journey on this phone. And I was genuinely shocked at how smooth the menu was. No lag whatsoever. This reminded me of a video that I did about the Galaxy S6, which is a phone from 10 years ago. And that phone aged like actual ass compared to this iPhone. Every menu lagged. And Google Maps was so slow that it would actually help you get more lost than you were before. But not the iPhone. Everything was just p- I lost like 20% in 10 minutes. What is going on? Maybe I could use something else to guide me to the right direction. Walt. This guy was watching me follow a pigeon for like a minute straight. As I continued walking, I managed to download a few of the apps, like Spotify and Google Keep. And despite the warning that I needed iOS 17 or later, I was simply able to download an older version, which, at least for now, works fine. But then what happens after? Do they just stop working? I was then looking at the back of the phone, and I was like, But then what happens after? Do they just stop working? I was then looking at the back of the phone, and I thought, There's no way that tiny camera is still good. The sensor is tiny. Surely it's gonna look like absolute dog, right? Well, I know that they definitely don't care. They're just pigeons. Anyway, if you go to the settings, you can actually change it to record at 4K, which is really surprising to me. What shocked me even more was when I actually tested the camera. This is a picture that I took on the S25. And this one was on the iPhone SE. Whoa. First of all, nowadays, the HDR effect is crazy. The sky did not look that blue. And everything just looks way too bright on the Samsung. Huh. When I look at Samsung's video, everything just looks too perfect and too colorful. But the iPhone, I don't know what it is, but you can just feel the aura. Look at it side by side. Doesn't the iPhone just feel more real? That is actually crazy. Time was running out. And so was the battery. I didn't want it to die before my first mission. That would be embarrassing. But then I kind of got lost. And whoa, the selfie camera is so shaky. And where even am I? Will I be late? Nevermind. We made it. Mission success. And I was a little sad to say goodbye to the original Snowball Mike. Farewell, buddy. Rest in peace. But we're not done just yet. Because I still have to make it home with less than 20%. And damn, that selfie camera looks so much better on Samsung. I'm not really surprised though. It's been a while since I had to double press the home button to launch the wallet. And paying for the train was a success. Am I the only one that runs up the stairs like that? And one last time when I tested the camera, I noticed that the phone does use HDR. But still, the photos come out looking more real. Why do the colors look so good? It actually reminds me of GTA San Andreas. The orange tint feels so nostalgic. But I was pretty concerned to see only 9% of my battery left after just one mission. Can I do anything else before I get home? Apparently, yes. But I accidentally logged into my wrong Instagram account. And for some reason, I was only following Eddie Hall. What is going on here? I'm not gonna lie. The scrolling was a little bit laggy. But since everything else was so good, I really cannot complain more about this phone. And what am I looking at? What is this real? Anyway, as good as the phone was, the battery died way too quickly before I could play any games or even listen to some music. Okay, so the phone lasted like two hours. That's a mission fail. But I'm not done yet. Ow. Of course, at this point, iPhones didn't even have wireless charging yet. But it's fine because there is one last thing that I want to do. I guess it's time I entertain myself in bed. I mean... Oh, and I almost forgot. I've got a headphone jack on this phone. So I can plug it in and actually listen to some premium music. Man, back then, music hit different when you had to spend like three hours trying to untangle your cables. That's childhood right there. I wish we got to own more things physically instead of having to stream everything. I even got a vinyl of my favorite album and I don't even own a vinyl player. Maybe it's a good thing that I have wireless headphones. I didn't realize how peaceful phones used to be. TikTok wasn't even a thing back then and people actually used to touch grass. Literally, while playing Pokemon Go. Back then, phones didn't have 200 megapixel cameras. You took a photo, it looked bad and you posted it anyway. I can't decide if back then I could feel more or I just knew less. I thought I missed the phone but in reality, I missed the time when I didn't care about anything."""
    fd = FeatureDetector()
    results = fd.detect_features(text)

    print(results)

    #print("Detected features:")
    #for feature, sentences in results.items():
        #print(f"\nFeature: {feature}")
        #for sentence in sentences:
            #print(f"  - {sentence}")
