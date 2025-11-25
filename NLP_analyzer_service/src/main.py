import signal
import time
import os
from NLP_service import SentimentAnalysisService

rabbit_url = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
service = SentimentAnalysisService(rabbit_url=rabbit_url)

def shutdown_handler(signum, frame):
    print(f"[*] Received signal {signum}, shutting down service...")
    service.close()
    exit(0)

# Register handlers
signal.signal(signal.SIGTERM, shutdown_handler)
signal.signal(signal.SIGINT, shutdown_handler)

# Start listening
service.start()
print("[*] NLP Analyzer service is running and listening to RabbitMQ...")

# Keep container alive
while True:
    time.sleep(10)
