# NLP_service.py
import json
import os
import pika
import sys
import threading
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../src")))

from consume import RabbitConsumer
from publish import RabbitPublisher
from feature_detection import FeatureDetector
from sentiment_analyzer import SentimentAnalyzer

class SentimentAnalysisService:
    """
    Main service that combines feature detection and sentiment analysis with publishing.
    """
    
    def __init__(self, rabbit_url=None):
        self.rabbit_url = rabbit_url or "amqp://guest:guest@localhost:5672/"
        self._running = False
        self._service_thread = None
        
        self.consumer = RabbitConsumer(rabbit_url)
        self.features = FeatureDetector()
        self.sentiments = SentimentAnalyzer()
        self.publisher = RabbitPublisher(rabbit_url)
        
    def analyze_and_publish(self):
        """Complete pipeline: analyze text and publish results."""
        print("[*] Starting background consumer...")
        self.consumer.start_background()
        self._running = True

        while self._running:
            print("[*] Waiting for next transcript message...")
            if not self.consumer.wait_for_message(timeout=5):  # Shorter timeout for responsiveness
                continue

            video_id = self.consumer.get_video_id()
            text = self.consumer.get_text()

            if not video_id or not text:
                print("[!] Received invalid or incomplete message, skipping...")
                continue

            print(f"[*] Analyzing text for video: {video_id}")
            try:
                feature_results = self.features.detect_features(text)
                print(f"[*] Detected features: {list(feature_results.keys())}")
                
                sentiment_results = self.sentiments.analyze_feature_sentiments(feature_results)
                print(f"[*] Generated sentiment results")
                
                self.publisher.publish_analysis(video_id, sentiment_results)
                print(f"[*] Published analysis for video {video_id}")
                
            except Exception as e:
                print(f"[*] Error processing video {video_id}: {e}")
                
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
                print("[*] Service thread did not finish gracefully")
        
        print("[*] Service shutdown complete")