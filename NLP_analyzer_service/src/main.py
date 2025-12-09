import signal
import time
import os
import json
import pika
from NLP_service import SentimentAnalysisService

DEFAULT_LOCAL = "amqp://guest:guest@localhost:5672/"
RABBITMQ_URL = os.getenv("AMQP_URL") or os.getenv("RABBITMQ_URL") or DEFAULT_LOCAL
EXCHANGE = "events"
ROUTING_KEY = "video.transcript.completed"

# Helper for safe RabbitMQ connection
def create_rabbit_connection(url, retries=5, delay=3):
    for attempt in range(1, retries + 1):
        try:
            connection = pika.BlockingConnection(pika.URLParameters(url))
            channel = connection.channel()
            # Ensure exchange exists
            channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
            return connection, channel
        except pika.exceptions.AMQPConnectionError as e:
            print(f"[!] RabbitMQ connection attempt {attempt} failed: {e}")
            time.sleep(delay)
    raise RuntimeError("Failed to connect to RabbitMQ after multiple attempts")

service = SentimentAnalysisService(rabbit_url=RABBITMQ_URL)

# Signal handling
def shutdown_handler(signum, frame):
    print(f"[*] Received signal {signum}, shutting down service...")
    service.close()
    exit(0)

signal.signal(signal.SIGTERM, shutdown_handler)
signal.signal(signal.SIGINT, shutdown_handler)

service.start()  # Assuming this launches the consumer in a background thread
print("[*] NLP Analyzer service is running and listening to RabbitMQ...")

while True:
    time.sleep(10)
