# consumer.py
import json
import pika

RABBITMQ_URL = "amqp://guest:guest@localhost:5672/"
EXCHANGE = "events"
QUEUE = "transcripts_queue"
ROUTING_KEY = "video.transcript.ready"


class RabbitConsumer:
    """Consumes transcript messages from RabbitMQ."""

    def __init__(self, rabbit_url=RABBITMQ_URL):
        self.rabbit_url = rabbit_url
        self.connection = None
        self.channel = None

    def connect(self):
        """Establish connection and declare exchange/queue."""
        params = pika.URLParameters(self.rabbit_url)
        self.connection = pika.BlockingConnection(params)
        self.channel = self.connection.channel()
        self.channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
        self.channel.queue_declare(queue=QUEUE, durable=True)
        self.channel.queue_bind(exchange=EXCHANGE, queue=QUEUE, routing_key=ROUTING_KEY)

    def on_message(self, ch, method, properties, body):
        """Callback for received messages."""
        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            print("[!] Received invalid JSON")
            ch.basic_ack(delivery_tag=method.delivery_tag)
            return

        print(f"[x] Received message: {json.dumps(data, indent=2)}")

        # For now, just acknowledge — later we’ll forward to analyzer
        ch.basic_ack(delivery_tag=method.delivery_tag)

    def start(self):
        """Start consuming messages."""
        self.connect()
        print("[*] Waiting for 'video.transcript.ready' messages...")
        self.channel.basic_consume(queue=QUEUE, on_message_callback=self.on_message)
        try:
            self.channel.start_consuming()
        except KeyboardInterrupt:
            print("\n[!] Stopping consumer...")
            self.stop()

    def stop(self):
        """Cleanly close connections."""
        if self.channel and self.channel.is_open:
            self.channel.stop_consuming()
        if self.connection and self.connection.is_open:
            self.connection.close()
        print("[x] Connection closed.")

