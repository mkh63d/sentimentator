# NLP_analyzer_service/src/publisher.py

import json
import os
import pika


class RabbitPublisher:
    def __init__(self, rabbit_url=None):
        self.rabbit_url = rabbit_url or os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
        self.params = pika.URLParameters(self.rabbit_url)
        self.connection = pika.BlockingConnection(self.params)
        self.channel = self.connection.channel()

        # Declare queue (or exchange if needed)
        self.queue_name = "video.analysis.complete"
        self.channel.queue_declare(queue=self.queue_name, durable=True)

    def publish_analysis(self, video_id: str, sentiment_results: dict):
        """
        Publish structured sentiment analysis result as JSON.
        """
        message = {
            "event": "video.analysis.complete",
            "video_id": video_id,
            "data": sentiment_results,
        }

        body = json.dumps(message)
        self.channel.basic_publish(
            exchange="",
            routing_key=self.queue_name,
            body=body,
            properties=pika.BasicProperties(
                delivery_mode=2,  # make message persistent
                content_type="application/json"
            ),
        )
        print(f"✅ Published sentiment analysis for {video_id} to {self.queue_name}")

    def close(self):
        self.connection.close()


if __name__ == "__main__":
    # quick example test
    sample_data = {
        "display": {"positive": 1, "negative": 0, "overall": "positive"},
        "battery": {"positive": 0, "negative": 1, "overall": "negative"}
    }

    pub = RabbitPublisher()
    pub.publish_analysis("video_001", sample_data)
    pub.close()
