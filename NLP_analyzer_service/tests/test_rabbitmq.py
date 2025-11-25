import pytest
import pika
import json
import time

RABBITMQ_URL = "amqp://guest:guest@localhost:5672/"
EXCHANGE = "test_exchange"
QUEUE = "test_queue"
ROUTING_KEY = "test.key"


@pytest.mark.integration
def test_rabbit_communication():
    # Connect to RabbitMQ
    params = pika.URLParameters(RABBITMQ_URL)
    connection = pika.BlockingConnection(params)
    channel = connection.channel()

    # Declare exchange and queue
    channel.exchange_declare(exchange=EXCHANGE, exchange_type='direct', durable=True)
    channel.queue_declare(queue=QUEUE, durable=True)
    channel.queue_bind(exchange=EXCHANGE, queue=QUEUE, routing_key=ROUTING_KEY)

    # Message to send
    test_message = {"analysis_id": 999, "transcript": "Integration test transcript."}

    # Publish message
    channel.basic_publish(
        exchange=EXCHANGE,
        routing_key=ROUTING_KEY,
        body=json.dumps(test_message),
        properties=pika.BasicProperties(delivery_mode=2),
    )

    # Consume message
    method_frame, header_frame, body = channel.basic_get(queue=QUEUE, auto_ack=True)
    
    # Retry briefly if queue is empty (RabbitMQ async delay)
    retries = 5
    while body is None and retries > 0:
        time.sleep(0.5)
        method_frame, header_frame, body = channel.basic_get(queue=QUEUE, auto_ack=True)
        retries -= 1

    assert body is not None, "No message received from RabbitMQ"

    received_message = json.loads(body)
    assert received_message == test_message, "Received message does not match sent message"

    connection.close()
