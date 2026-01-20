import os
import json
import redis
import time
from src.tasks import process_question
from dotenv import load_dotenv

load_dotenv()

# Configuration
REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost:6379')
QUEUE_NAME = 'inference_queue'

# Connect to Redis
try:
    r = redis.from_url(REDIS_URL)
    # Ping to verify connection
    r.ping()
    print(f"Connected to Redis at {REDIS_URL}")
except redis.ConnectionError as e:
    print(f"Failed to connect to Redis: {e}")
    exit(1)

def handle_task(task_data):
    """
    Dispatch the task to the appropriate function.
    """
    task_type = task_data.get('type')
    job_id = task_data.get('job_id')
    
    print(f"Received task: {task_type} (Job ID: {job_id})")

    result = None
    
    if task_type == 'process_question':
        # Extract arguments
        question = task_data.get('payload', {}).get('question')
        user_id = task_data.get('payload', {}).get('user_id')
        
        if question and user_id:
            # CALL THE FUNCTION
            result = process_question(question, user_id)
        else:
            result = {"error": "Missing question or user_id", "status": "failed"}
            
    else:
        print(f"Unknown task type: {task_type}")
        result = {"error": f"Unknown task type: {task_type}", "status": "failed"}

    # Save the result back to Redis for the backend to retrieve
    # Key format: "job_result:{job_id}"
    if job_id:
        result_key = f"job_result:{job_id}"
        # Expire result after 1 hour to save memory
        r.setex(result_key, 3600, json.dumps(result))
        print(f"Result saved to {result_key}")

def start_worker():
    print(f"Worker listening on queue: '{QUEUE_NAME}'...")
    
    while True:
        try:
            # BLPOP blocks until an item is available in the list
            # It returns a tuple: (queue_name, data)
            _, message = r.blpop(QUEUE_NAME)
            
            # Parse JSON
            task_data = json.loads(message)
            
            # Process
            handle_task(task_data)
            
        except json.JSONDecodeError:
            print("Error: Failed to decode JSON message")
        except Exception as e:
            print(f"Error in worker loop: {e}")
            # Sleep briefly to avoid tight loop in case of persistent errors
            time.sleep(1)

if __name__ == '__main__':
    start_worker()
