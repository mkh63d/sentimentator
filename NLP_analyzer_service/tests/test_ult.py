import os
import sys
import json
import time
import pytest
import threading
import pika

# run this with pytest -v -s test_ult.py from cmd

# Add src folder to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../src")))
from NLP_service import SentimentAnalysisService

RABBITMQ_URL = "amqp://guest:guest@localhost:5672/"
EXCHANGE = "events"
SENDING_ROUTING_KEY = "video.transcript.ready"
RESULT_QUEUE = "video.analysis.complete"

RESOURCE_DIR = os.path.join(os.path.dirname(__file__), "../resources")
MOCK_TRANSCRIPT_PATH = os.path.join(RESOURCE_DIR, "review_text.txt")

EXPECTED_ANALYSIS = {
    "camera": {'positive_score': 0.9202831983566284, 'negative_score': 0.0, 'neutral_score': 0.44070518016815186, 'overall': 'positive'},
    "battery": {'positive_score': 0.5078599452972412, 'negative_score': 0.5051186680793762, 'neutral_score': 0.0, 'overall': 'neutral'},
    "display": {'positive_score': 0.7228912711143494, 'negative_score': 0.0, 'neutral_score': 0.0, 'overall': 'positive'},
    "speaker": {'positive_score': 0.0, 'negative_score': 0.5635719299316406, 'neutral_score': 0.0, 'overall': 'negative'},
    "performance": {'positive_score': 0.6023634076118469, 'negative_score': 0.41846367716789246, 'neutral_score': 0.0, 'overall': 'positive'}
}

# ---------------- Helpers ----------------
class MockSender:
    """Sends a transcript to RabbitMQ like a video event."""
    def __init__(self, video_id=1):
        self.video_id = video_id

    def send(self):
        with open(MOCK_TRANSCRIPT_PATH, "r", encoding="utf-8") as f:
            transcript = f.read()
        
        # Use the format expected by the consumer
        message = {
            "video_id": self.video_id, 
            "transcript": transcript,
            "event": "video.transcript.ready"  # Add event type
        }

        connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
        channel = connection.channel()
        channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
        channel.basic_publish(
            exchange=EXCHANGE,
            routing_key=SENDING_ROUTING_KEY,
            body=json.dumps(message),
            properties=pika.BasicProperties(
                content_type="application/json",
                delivery_mode=2  # persistent
            )
        )
        connection.close()
        print(f"[*] Sent message for video {self.video_id}")
        return message


class MockReceiver:
    """Listens to the analysis queue and returns the processed result."""
    def receive(self, timeout=15):
        connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
        channel = connection.channel()
        channel.queue_declare(queue=RESULT_QUEUE, durable=True)

        result = None
        method_frame = None
        properties = None
        body = None
        
        # Get one message
        method_frame, properties, body = channel.basic_get(queue=RESULT_QUEUE)
        if method_frame:
            result = json.loads(body)
            channel.basic_ack(method_frame.delivery_tag)
            print(f"[*] Received analysis result for video {result.get('video_id', 'unknown')}")
        else:
            print("[*] No message received from analysis queue")
            
        channel.close()
        connection.close()
        return result

# ---------------- Test ----------------
@pytest.mark.integration
def test_ultimate_sentiment_pipeline():
    # Start the service in a background thread
    service = SentimentAnalysisService(rabbit_url=RABBITMQ_URL)
    
    # FIX: Pass the function reference, don't call it immediately
    thread = service.start()  # Use the service's start method
    
    time.sleep(3)  # Give service more time to connect and set up

    try:
        # Send the transcript
        sender = MockSender(video_id=1)
        sent_message = sender.send()
        video_id = sent_message["video_id"]

        # Wait a bit for processing
        time.sleep(5)

        # Receive the processed analysis
        receiver = MockReceiver()
        received = receiver.receive(timeout=10)

        # --- Assertions ---
        assert received is not None, "No message received from service"
        assert received.get("video_id") == video_id, f"Video ID mismatch: {received.get('video_id')} != {video_id}"
        
        data = received.get("data", {})
        assert data, "No data in received message"

        # Check if we got the expected features
        found_features = list(data.keys())
        print(f"[*] Found features: {found_features}")

        for feature, expected_scores in EXPECTED_ANALYSIS.items():
            if feature in data:
                actual_scores = data[feature]
                for key in expected_scores:
                    if key in actual_scores:
                        if isinstance(expected_scores[key], float):
                            diff = abs(actual_scores[key] - expected_scores[key])
                            assert diff < 0.15, f"Mismatch for {feature}:{key} - expected {expected_scores[key]}, got {actual_scores[key]} (diff: {diff})"
                        else:
                            assert actual_scores[key] == expected_scores[key], f"Mismatch for {feature}:{key} - expected {expected_scores[key]}, got {actual_scores[key]}"
                    else:
                        print(f"[*]  Key {key} not found in {feature} data")
            else:
                print(f"[*]  Feature {feature} not found in results")

        print("✓ All assertions passed!")

    except Exception as e:
        print(f"X Test failed: {e}")
        raise
    finally:
        # Clean up - use service.close() which handles everything
        print("[*] Cleaning up...")
        service.close()