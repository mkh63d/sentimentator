# consume.py
import json
import pika
import threading
import time
import os

DEFAULT_LOCAL = "amqp://guest:guest@localhost:5672/"
RABBITMQ_URL = os.getenv("RABBITMQ_URL", DEFAULT_LOCAL)
EXCHANGE = "sentimentator.topic"
QUEUE = "video.transcript"
ROUTING_KEY = "video.transcript.completed"



class RabbitConsumer:
    """Consumes transcript messages from RabbitMQ and stores the latest one."""

    def __init__(self, rabbit_url=RABBITMQ_URL):
        self.rabbit_url = self._clean_url(rabbit_url)
        self.connection = None
        self.channel = None
        self._latest_message = None
        self._lock = threading.Lock()
        self._new_message_event = threading.Event()
        self._consuming = False
        self._consumer_thread = None

    def _clean_url(self, url):
        """Clean up URL to handle double slashes (e.g., ending with //)."""
        # If URL ends with //, replace with single /
        if url.endswith('//'):
            url = url[:-1]
        return url

    def connect(self):
        """Establish connection with error handling"""
        try:
            params = pika.URLParameters(self.rabbit_url)
            self.connection = pika.BlockingConnection(params)
            self.channel = self.connection.channel()
            self.channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
            self.channel.queue_declare(queue=QUEUE, durable=True)
            self.channel.queue_bind(exchange=EXCHANGE, queue=QUEUE, routing_key=ROUTING_KEY)
            return True
        except Exception as e:
            print(f"[!] Failed to connect to RabbitMQ: {e}")
            return False

    def _on_message(self, ch, method, properties, body):
        try:
            data = json.loads(body)
            print(f"[*] Received message for video: {data.get('analysis_id', 'unknown')}")

            with self._lock:
                self._latest_message = data
                self._new_message_event.set()

            ch.basic_ack(delivery_tag=method.delivery_tag)

        except json.JSONDecodeError:
            print("[!] Received invalid JSON")
            ch.basic_ack(delivery_tag=method.delivery_tag)
        except Exception as e:
            print(f"[!] Error processing message: {e}")
            ch.basic_ack(delivery_tag=method.delivery_tag)

    def start_background(self):
        """Start consuming in a background thread."""
        if not self.connect():
            return False

        def _consume():
            print("[*] Waiting for 'video.transcript.completed' messages...")
            try:
                self.channel.basic_consume(queue=QUEUE, on_message_callback=self._on_message)
                self._consuming = True
                self.channel.start_consuming()
            except Exception as e:
                print(f"[!] Consumer stopped: {e}")
            finally:
                self._consuming = False

        self._consumer_thread = threading.Thread(target=_consume, daemon=True)
        self._consumer_thread.start()
        return True

    def wait_for_message(self, timeout=None):
        """Block until a new message arrives or timeout."""
        received = self._new_message_event.wait(timeout)
        if received:
            self._new_message_event.clear()
        return received

    def get_latest_message(self):
        """Return the full latest message dict (thread-safe)."""
        with self._lock:
            return dict(self._latest_message) if self._latest_message else None

    # optional: keep these but implemented from latest_message
    def get_analysis_id(self):
        msg = self.get_latest_message()
        return msg.get("analysis_id") if msg else None

    def get_text(self):
        msg = self.get_latest_message()
        # keep backward compatibility: some messages may nest transcript under 'data'
        return msg.get("transcript") or (msg.get("data") or {}).get("transcript") if msg else None

    def get_video_id(self):
        msg = self.get_latest_message()
        return msg.get("video_id") if msg else None

    def get_original_url(self):
        msg = self.get_latest_message()
        return msg.get("original_url") if msg else None

    def stop(self):
        """Safely stop the consumer."""
        print("[*] Stopping RabbitMQ consumer...")

        # Stop consuming first
        if self.channel and self.channel.is_open and self._consuming:
            try:
                self.channel.stop_consuming()
                print("[*] Stopped consuming messages")
            except Exception as e:
                print(f"[*] Error stopping consumption: {e}")

        # Close channel
        if self.channel and self.channel.is_open:
            try:
                self.channel.close()
                print("[*] Closed channel")
            except Exception as e:
                print(f"[*] Error closing channel: {e}")

        # Close connection
        if self.connection and self.connection.is_open:
            try:
                self.connection.close()
                print("[*] Closed connection")
            except Exception as e:
                print(f"[*] Error closing connection: {e}")

        self._consuming = False
        print("[*] RabbitMQ consumer stopped.")
