# publisher.py
import json
import os
import pika


class RabbitPublisher:
    def __init__(self, rabbit_url=None):
        self.rabbit_url = rabbit_url or os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
        self.connection = None
        self.channel = None
        self._connect()

    def _connect(self):
        """Establish connection with error handling"""
        try:
            self.params = pika.URLParameters(self.rabbit_url)
            self.connection = pika.BlockingConnection(self.params)
            self.channel = self.connection.channel()
            self.queue_name = "video.analysis.complete"
            self.channel.queue_declare(queue=self.queue_name, durable=True)
            print("[✓] RabbitMQ publisher connected")
        except Exception as e:
            print(f"[!] Failed to connect RabbitMQ publisher: {e}")
            raise

    def publish_analysis(self, video_id: str, sentiment_results: dict):
        """
        Publish structured sentiment analysis result as JSON.
        """
        # Reconnect if connection is lost
        if not self.connection or self.connection.is_closed:
            print("[!] Publisher connection lost, reconnecting...")
            self._connect()

        message = {
            "event": "video.analysis.complete",
            "video_id": video_id,
            "data": sentiment_results,
        }

        try:
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
            print(f"[*] Published sentiment analysis for {video_id} to {self.queue_name}")
        except Exception as e:
            print(f"[*] Failed to publish message: {e}")
            # Try to reconnect and retry once
            try:
                self._connect()
                self.channel.basic_publish(
                    exchange="",
                    routing_key=self.queue_name,
                    body=body,
                    properties=pika.BasicProperties(
                        delivery_mode=2,
                        content_type="application/json"
                    ),
                )
                print(f"[*] Published sentiment analysis for {video_id} (after reconnect)")
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