import os
import time
from langchain_community.chat_models import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser

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
