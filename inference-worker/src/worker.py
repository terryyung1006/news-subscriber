import os
import json
import time
import pika
from src.tasks import process_question
from dotenv import load_dotenv

load_dotenv()

RABBITMQ_URL = os.getenv('RABBITMQ_URL', 'amqp://guest:guest@localhost:5672/')
REQUESTS_QUEUE = 'inference.requests'
RESULTS_QUEUE = 'inference.results'


def handle_delivery(channel, method, properties, body):
    try:
        data = json.loads(body)
    except json.JSONDecodeError:
        print("Error: Failed to decode JSON message")
        channel.basic_ack(delivery_tag=method.delivery_tag)
        return

    correlation_id = data.get('correlation_id', '')
    user_id = data.get('user_id', '')
    user_name = data.get('user_name', 'User')
    question = data.get('question', '')

    print(f"Received task (correlation_id={correlation_id})")

    if not question or not user_id:
        result = {
            'correlation_id': correlation_id,
            'answer': '',
            'status': 'failed',
            'error': 'Missing question or user_id',
        }
    else:
        try:
            raw = process_question(question, user_id, user_name)
            result = {
                'correlation_id': correlation_id,
                'answer': raw.get('response', ''),
                'status': raw.get('status', 'completed'),
                'error': raw.get('error', ''),
            }
        except Exception as e:
            result = {
                'correlation_id': correlation_id,
                'answer': '',
                'status': 'failed',
                'error': str(e),
            }

    channel.basic_publish(
        exchange='',
        routing_key=RESULTS_QUEUE,
        body=json.dumps(result),
        properties=pika.BasicProperties(
            delivery_mode=2,
            content_type='application/json',
            correlation_id=correlation_id,
        ),
    )

    channel.basic_ack(delivery_tag=method.delivery_tag)
    print(f"Result published for correlation_id={correlation_id}")


def start_worker():
    while True:
        try:
            params = pika.URLParameters(RABBITMQ_URL)
            params.heartbeat = 600
            connection = pika.BlockingConnection(params)
            channel = connection.channel()

            channel.queue_declare(queue=REQUESTS_QUEUE, durable=True)
            channel.queue_declare(queue=RESULTS_QUEUE, durable=True)
            channel.basic_qos(prefetch_count=1)

            channel.basic_consume(
                queue=REQUESTS_QUEUE,
                on_message_callback=handle_delivery,
            )

            print(f"Worker listening on '{REQUESTS_QUEUE}'...")
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            print(f"Connection error: {e}. Retrying in 5s...")
            time.sleep(5)
        except KeyboardInterrupt:
            print("Worker stopped.")
            break
        except Exception as e:
            print(f"Unexpected error: {e}. Retrying in 5s...")
            time.sleep(5)


if __name__ == '__main__':
    start_worker()
