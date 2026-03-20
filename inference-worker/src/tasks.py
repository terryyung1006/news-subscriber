import os
import json
import re
from langchain_community.chat_models import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_core.messages import SystemMessage, HumanMessage, AIMessage

# Initialize the LLM (Ollama)
# Make sure the model is pulled: `docker exec -it news_subscriber_ollama ollama pull deepseek-r1:1.5b`
llm = ChatOllama(
    model="deepseek-r1:1.5b",
    base_url="http://localhost:11434",
    temperature=0.7,
)

def process_question(question: str, user_id: str, user_name: str = "User"):
    """
    Task to process a user question using an LLM (Ollama).
    """
    print(f"Processing question for user {user_name} ({user_id}): {question}")
    
    try:
        prompt = ChatPromptTemplate.from_template(
            "You are a helpful assistant. The user's name is {user_name}.\n"
            "User Question: {question}"
        )
        
        chain = prompt | llm | StrOutputParser()
        
        # Invoke the chain
        response = chain.invoke({"user_name": user_name, "question": question})
        
        print(f"Generated response: {response}")
        return {
            "user_id": user_id,
            "user_name": user_name,
            "question": question,
            "response": response,
            "status": "completed"
        }
    except Exception as e:
        print(f"Error processing question: {e}")
        return {
            "user_id": user_id,
            "question": question,
            "error": str(e),
            "status": "failed"
        }


ONBOARDING_SYSTEM_PROMPT = """You are a friendly AI news assistant helping a new user set up their preferences. Your personality is like a personal robot that just woke up and is eager to learn about its owner.

Your goals:
1. Ask thoughtful follow-up questions to understand the user's news interests
2. Extract specific, actionable preferences that can be used to personalize news
3. After gathering enough information (at least 3 distinct preferences), extract and summarize them

Guidelines:
- Be warm and conversational, staying in the "personal robot" character
- Ask ONE follow-up question at a time to dig deeper
- Look for: industries, topics, regions, specific companies, reading habits
- After 2-3 exchanges, if you have enough info, complete the onboarding

IMPORTANT: You must respond with valid JSON only. No other text before or after the JSON.

If you need to ask more questions, respond with:
{"is_complete": false, "memories": [], "response": "Your follow-up question here..."}

When you have gathered at least 3 distinct preferences, respond with:
{"is_complete": true, "memories": [{"title": "Short Title", "content": "Detailed preference description", "category": "news_preference"}, ...], "response": "Your completion message summarizing what you learned..."}

Categories to use: "news_preference", "interest", "reading_habit"
"""


def process_onboarding(conversation: list, user_id: str, user_name: str = "User"):
    """
    Process onboarding conversation and extract memories when ready.

    Args:
        conversation: List of message dicts with 'role' and 'content'
        user_id: The user's ID
        user_name: The user's name for personalization

    Returns:
        Dict with response, is_complete flag, and extracted memories
    """
    print(f"Processing onboarding for user {user_name} ({user_id})")
    print(f"Conversation length: {len(conversation)} messages")

    try:
        # Build messages for LLM
        messages = [SystemMessage(content=ONBOARDING_SYSTEM_PROMPT.replace("{user_name}", user_name))]

        for msg in conversation:
            if msg.get('role') == 'assistant':
                messages.append(AIMessage(content=msg.get('content', '')))
            elif msg.get('role') == 'user':
                messages.append(HumanMessage(content=msg.get('content', '')))

        # Call LLM
        response = llm.invoke(messages)
        response_text = response.content.strip()

        print(f"Raw LLM response: {response_text}")

        # Parse JSON from response - handle potential markdown code blocks
        json_text = response_text
        if "```json" in response_text:
            json_match = re.search(r'```json\s*(.*?)\s*```', response_text, re.DOTALL)
            if json_match:
                json_text = json_match.group(1)
        elif "```" in response_text:
            json_match = re.search(r'```\s*(.*?)\s*```', response_text, re.DOTALL)
            if json_match:
                json_text = json_match.group(1)

        # Try to find JSON object in the response
        json_match = re.search(r'\{.*\}', json_text, re.DOTALL)
        if json_match:
            json_text = json_match.group(0)

        try:
            result = json.loads(json_text)
        except json.JSONDecodeError:
            # If JSON parsing fails, treat as a regular response
            print(f"Failed to parse JSON, treating as regular response")
            result = {
                "is_complete": False,
                "memories": [],
                "response": response_text
            }

        # Validate structure
        is_complete = result.get("is_complete", False)
        memories = result.get("memories", [])
        response_msg = result.get("response", "")

        # Ensure we have at least 3 memories if marked complete
        if is_complete and len(memories) < 3:
            is_complete = False
            response_msg = "I'd love to learn more about you! " + response_msg

        return {
            "response": response_msg,
            "is_complete": is_complete,
            "memories": memories,
            "status": "completed"
        }

    except Exception as e:
        print(f"Error processing onboarding: {e}")
        return {
            "response": "",
            "is_complete": False,
            "memories": [],
            "error": str(e),
            "status": "failed"
        }
