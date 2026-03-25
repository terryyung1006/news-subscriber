import os
import json
import time
import threading
import pika
import redis
from src.tasks import process_question, process_onboarding
from dotenv import load_dotenv

load_dotenv()

RABBITMQ_URL = os.getenv('RABBITMQ_URL', 'amqp://guest:guest@localhost:5672/')
REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost:6379')
REQUESTS_QUEUE = 'inference.requests'
RESULTS_QUEUE = 'inference.results'
REDIS_INFERENCE_QUEUE = 'inference_queue'


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


def start_rabbitmq_worker():
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

            print(f"RabbitMQ worker listening on '{REQUESTS_QUEUE}'...")
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            print(f"RabbitMQ connection error: {e}. Retrying in 5s...")
            time.sleep(5)
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"RabbitMQ unexpected error: {e}. Retrying in 5s...")
            time.sleep(5)


def start_redis_worker():
    r = redis.from_url(REDIS_URL)
    print(f"Redis worker listening on '{REDIS_INFERENCE_QUEUE}'...")

    while True:
        try:
            # BLPOP blocks until a message is available
            item = r.blpop(REDIS_INFERENCE_QUEUE, timeout=5)
            if item is None:
                continue

            _, raw_data = item
            data = json.loads(raw_data)

            job_id = data.get('job_id', '')
            task_type = data.get('type', '')
            payload = data.get('payload', {})

            print(f"Redis task received (job_id={job_id}, type={task_type})")

            if task_type == 'process_onboarding':
                user_id = payload.get('user_id', '')
                user_name = payload.get('user_name', 'User')
                conversation = payload.get('conversation', [])

                try:
                    raw = process_onboarding(conversation, user_id, user_name)
                    result = {
                        'response': raw.get('response', ''),
                        'is_complete': raw.get('is_complete', False),
                        'memories': raw.get('memories', []),
                        'status': raw.get('status', 'completed'),
                        'error': raw.get('error', ''),
                    }
                except Exception as e:
                    result = {
                        'response': '',
                        'is_complete': False,
                        'memories': [],
                        'status': 'failed',
                        'error': str(e),
                    }
            else:
                result = {
                    'status': 'failed',
                    'error': f'Unknown task type: {task_type}',
                }

            result_key = f"job_result:{job_id}"
            r.set(result_key, json.dumps(result), ex=300)
            print(f"Redis result published for job_id={job_id}")

        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"Redis worker error: {e}. Retrying in 5s...")
            time.sleep(5)


if __name__ == '__main__':
    # Run both workers in parallel
    rabbitmq_thread = threading.Thread(target=start_rabbitmq_worker, daemon=True)
    redis_thread = threading.Thread(target=start_redis_worker, daemon=True)

    rabbitmq_thread.start()
    redis_thread.start()

    print("Worker started (RabbitMQ + Redis)")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("Worker stopped.")
