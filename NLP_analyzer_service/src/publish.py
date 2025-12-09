# publisher.py
import json
import os
import pika

DEFAULT_LOCAL = "amqp://guest:guest@localhost:5672/"
RABBITMQ_URL = os.getenv("RABBITMQ_URL", DEFAULT_LOCAL)

class RabbitPublisher:
    def __init__(self, rabbit_url=None):
        self.rabbit_url = self._clean_url(rabbit_url or RABBITMQ_URL)
        self.connection = None
        self.channel = None
        self._connect()

    def _clean_url(self, url):
        """Clean up URL to handle double slashes (e.g., ending with //)."""
        # If URL ends with //, replace with single /
        if url.endswith('//'):
            url = url[:-1]
        return url

    def _connect(self):
        try:
            self.params = pika.URLParameters(self.rabbit_url)
            self.connection = pika.BlockingConnection(self.params)
            self.channel = self.connection.channel()
            self.queue_name = "video.analysis"
            self.channel.exchange_declare(exchange="sentimentator.topic", exchange_type="topic", durable=True)
            self.channel.queue_declare(queue=self.queue_name, durable=True)
            self.channel.queue_bind(
                queue=self.queue_name,
                exchange="sentimentator.topic",
                routing_key="video.analysis.completed"
            )
            print("[*] RabbitMQ publisher connected")
        except Exception as e:
            print(f"[!] Failed to connect RabbitMQ publisher: {e}")
            raise


    def publish_analysis(self,
                         analysis_id: str,
                         original_url: str,
                         video_id: str,
                         transcript: str,
                         sentiment_results: dict):
        """
        Publish AnalysisCompletedEvent:
        {
          analysis_id, original_url, video_id, transcript,
          sentiment: { overall, features, confidence }
        }
        """
        if not self.connection or self.connection.is_closed:
            print("[!] Publisher connection lost, reconnecting...")
            self._connect()

        # Ensure sentiment_results is already in the SentimentAnalysis shape:
        # { "overall": float, "features": {...}, "confidence": float }
        message = {
            "event": "video.analysis.completed",
            "analysis_id": analysis_id,
            "original_url": original_url,
            "video_id": video_id,
            "transcript": transcript,
            "sentiment": sentiment_results,
        }

        try:
            body = json.dumps(message)
            # publish to exchange or direct queue depending on topology
            # here: publish to EXCHANGE with routing key for subscribers
            self.channel.basic_publish(
                exchange="sentimentator.topic",
                routing_key="video.analysis.completed",
                body=body,
                properties=pika.BasicProperties(
                    delivery_mode=2,
                    content_type="application/json"
                ),
            )
            print(f"[*] Published sentiment analysis for {analysis_id}")
        except Exception as e:
            print(f"[*] Failed to publish message: {e}")
            # reconnect and retry once (same as before)
            try:
                self._connect()
                self.channel.basic_publish(
                    exchange="sentimentator.topic",
                    routing_key="video.analysis.completed",
                    body=body,
                    properties=pika.BasicProperties(
                        delivery_mode=2,
                        content_type="application/json"
                    ),
                )
                print(f"[*] Published sentiment analysis for {analysis_id} (after reconnect)")
            except Exception as retry_error:
                print(f"[*] Failed to publish even after reconnect: {retry_error}")
                raise

    def close(self):
        """Safely close the publisher."""
        print("[*] Closing RabbitMQ publisher...")
        try:
            if self.channel and self.channel.is_open:
                self.channel.close()
            if self.connection and self.connection.is_open:
                self.connection.close()
            print("[*] RabbitMQ publisher closed")
        except Exception as e:
            print(f"[*] Error closing publisher (ignoring): {e}")
