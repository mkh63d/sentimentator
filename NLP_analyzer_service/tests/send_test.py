import pika
import json

RABBITMQ_URL = RABBITMQ_URL = "amqp://guest:guest@sentimentator-rabbitmq-1:5672/"
EXCHANGE = "events"
ROUTING_KEY = "video.transcript.completed"

# Sample message
message = {
    "analysis_id": 1,
    "transcript": "This is a mock transcript for testing.",
    "event": "video.transcript.completed"
}

connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
channel = connection.channel()

# Make sure the exchange exists
channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)

# Publish the message
channel.basic_publish(
    exchange=EXCHANGE,
    routing_key=ROUTING_KEY,
    body=json.dumps(message),
    properties=pika.BasicProperties(
        content_type="application/json",
        delivery_mode=2
    )
)

connection.close()
print("[*] Message sent!")
