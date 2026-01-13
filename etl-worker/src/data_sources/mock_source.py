"""
Mock financial news data source for testing and development.

This implementation provides realistic mock data for testing the RAG system
without requiring external API connections. Now enhanced with LangChain + Ollama LLM generation.
"""

import random
import uuid
from datetime import datetime, timedelta
from typing import List, Optional

from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from langchain_community.llms import Ollama
from loguru import logger

from .base import DataSource, NewsItem
from ..load_config import get_config


class MockFinancialNewsSource(DataSource):
    """
    Mock data source that generates realistic financial news for testing.

    This class follows the Liskov Substitution Principle by properly
    implementing the DataSource interface.
    """

    def __init__(
        self, news_count: int = 10, use_llm: bool = True, model_name: str = "llama2"
    ):
        """
        Initialize the mock data source.

        Args:
            news_count: Number of mock news items to generate
            use_llm: Whether to use local LLM for content generation
            model_name: Name of the Ollama model to use
        """
        self.news_count = news_count
        self.use_llm = use_llm
        self.model_name = model_name
        self._connected = False
        self._mock_news = []
        self._llm_chain = None

        if self.use_llm:
            self._initialize_llm()

        self._generate_mock_data()

    def _initialize_llm(self) -> None:
        """Initialize the LangChain LLM with Ollama."""
        try:
            config = get_config()
            # Initialize Ollama LLM
            self._llm = Ollama(
                model=self.model_name,
                base_url=config.get_string("OLLAMA_BASE_URL", "http://localhost:11434"),
                temperature=config.get("OLLAMA_TEMPERATURE", 0.8),
                timeout=config.get_int("OLLAMA_TIMEOUT", 30),
            )

            # Create a prompt template for financial news generation
            prompt_template = PromptTemplate(
                input_variables=["category", "sentiment", "source"],
                template="""
You are a financial news journalist. Write a {sentiment} news article about {category} that would appear in {source}.

Requirements:
- Write a catchy, professional headline (max 80 characters)
- Write 2-3 paragraphs of realistic financial news content (around 70 words)
- Keep it factual and professional
- Include specific numbers, percentages, or financial data
- Make it sound like real financial news

Format your response exactly like this:
HEADLINE: [Your headline here] (make sure this is one line)
CONTENT: [Your content here] (make sure this is one line)
""",
            )

            # Create the LLM chain
            self._llm_chain = LLMChain(llm=self._llm, prompt=prompt_template)

            logger.info(
                f"Initialized LangChain LLM with Ollama model: {self.model_name}"
            )

        except ImportError as e:
            logger.warning(f"Could not import LangChain: {e}")
            logger.warning("Falling back to hardcoded mock data")
            self.use_llm = False
        except Exception as e:
            logger.warning(f"Error initializing LLM: {e}")
            logger.warning("Falling back to hardcoded mock data")
            self.use_llm = False

    def _generate_llm_content(
        self, category: str, sentiment: str, source: str
    ) -> tuple[str, str]:
        """
        Generate news title and content using LangChain + Ollama.

        Args:
            category: Financial category (stocks, bonds, etc.)
            sentiment: News sentiment (positive, negative, neutral)
            source: News source name

        Returns:
            tuple: (title, content)
        """
        if not self._llm_chain:
            return self._generate_fallback_content(category, sentiment)

        try:
            # Generate content using LangChain
            result = self._llm_chain.invoke(
                {"category": category, "sentiment": sentiment, "source": source}
            )

            # Parse the generated text
            if result:
                # Handle the new response format from invoke
                if "text" in result:
                    result_text = result["text"]
                else:
                    result_text = str(result)

                lines = result_text.strip().split("\n")
                title = ""
                content = ""

                for line in lines:
                    line = line.strip()
                    if line.startswith("HEADLINE:"):
                        title = line.replace("HEADLINE:", "").strip()
                    elif line.startswith("CONTENT:"):
                        content = line.replace("CONTENT:", "").strip()

                if title and content:
                    # Clean up the content
                    title = title.strip('"').strip()
                    content = content.strip()
                    return title, content

            # Fallback if parsing fails
            logger.warning("Failed to parse LLM response, using fallback")
            return self._generate_fallback_content(category, sentiment)

        except Exception as e:
            logger.warning(f"LLM generation failed: {e}")
            return self._generate_fallback_content(category, sentiment)

    def _generate_fallback_content(
        self, category: str, sentiment: str
    ) -> tuple[str, str]:
        """Generate fallback content when LLM fails."""
        # Enhanced fallback templates
        templates = {
            "stocks": {
                "positive": (
                    "Tech Stocks Rally on Strong Earnings",
                    "Technology stocks showed strong performance today as major companies reported better-than-expected earnings. The S&P 500 technology sector gained 2.3%, led by strong quarterly results from major tech firms.",
                ),
                "negative": (
                    "Market Selloff Continues Amid Economic Concerns",
                    "Stock markets experienced another day of losses as investors remain cautious about economic outlook. The Dow Jones Industrial Average fell 1.2% while the NASDAQ Composite dropped 1.8%.",
                ),
                "neutral": (
                    "Markets Trade Sideways Awaiting Economic Data",
                    "Stock markets showed little movement today as traders await key economic data releases. The major indices remained within 0.5% of their opening levels throughout the session.",
                ),
            },
            "bonds": {
                "positive": (
                    "Bond Yields Decline as Investors Seek Safety",
                    "Government bond yields fell today as investors seek safe-haven assets amid market uncertainty. The 10-year Treasury yield dropped 8 basis points to 4.15%.",
                ),
                "negative": (
                    "Bond Market Volatility Increases",
                    "Bond markets experienced increased volatility amid uncertainty about monetary policy. The 10-year Treasury yield swung between 4.20% and 4.35% during the session.",
                ),
                "neutral": (
                    "Bond Markets Stable with Minimal Movement",
                    "Bond markets remained relatively stable today with minimal price movements. The 10-year Treasury yield closed unchanged at 4.25%.",
                ),
            },
            "crypto": {
                "positive": (
                    "Cryptocurrency Surge Continues",
                    "Digital assets gained significant value today as institutional adoption increases. Bitcoin rose 5.2% to $67,500 while Ethereum gained 4.8% to $3,200.",
                ),
                "negative": (
                    "Crypto Market Correction Intensifies",
                    "Cryptocurrency markets experienced a correction today after recent gains. Bitcoin fell 3.1% to $62,800 while Ethereum declined 2.8% to $3,050.",
                ),
                "neutral": (
                    "Crypto Markets Consolidate",
                    "Cryptocurrency markets showed consolidation today with mixed performance. Bitcoin traded sideways around $65,000 while altcoins showed varied results.",
                ),
            },
            "forex": {
                "positive": (
                    "US Dollar Strengthens Against Major Currencies",
                    "The US dollar gained ground against major currencies today, with the dollar index rising 0.8%. The euro fell 0.6% to $1.0850 while the yen weakened to 150.20.",
                ),
                "negative": (
                    "Dollar Weakens on Economic Data",
                    "The US dollar weakened against major currencies following disappointing economic data. The dollar index fell 0.5% as the euro rose to $1.0950.",
                ),
                "neutral": (
                    "Forex Markets Show Limited Movement",
                    "Foreign exchange markets showed limited movement today with most major currency pairs trading within narrow ranges.",
                ),
            },
            "commodities": {
                "positive": (
                    "Oil Prices Rally on Supply Concerns",
                    "Oil prices rose today as supply concerns outweigh demand worries. Brent crude gained 2.1% to $82.50 per barrel while WTI rose 1.8% to $78.20.",
                ),
                "negative": (
                    "Commodity Prices Decline",
                    "Commodity prices fell across the board today amid concerns about global economic growth. Gold dropped 1.2% to $1,950 per ounce while silver fell 1.8%.",
                ),
                "neutral": (
                    "Commodity Markets Mixed",
                    "Commodity markets showed mixed performance today with energy prices rising while precious metals declined slightly.",
                ),
            },
            "economy": {
                "positive": (
                    "Economic Data Shows Recovery Signs",
                    "Recent economic data indicates signs of recovery in key sectors. GDP growth for the quarter was revised upward to 2.8% while unemployment claims fell to a 10-month low.",
                ),
                "negative": (
                    "Economic Indicators Point to Slowdown",
                    "Economic indicators suggest a potential slowdown in economic activity. Consumer confidence fell to its lowest level in six months while retail sales declined unexpectedly.",
                ),
                "neutral": (
                    "Economic Data Mixed",
                    "Economic data released today showed mixed results with some indicators improving while others remained unchanged from previous readings.",
                ),
            },
        }

        # Get template or use generic
        if category in templates and sentiment in templates[category]:
            return templates[category][sentiment]

        # Generic fallback
        generic_templates = {
            "positive": (
                "Positive Market News",
                "Financial markets showed positive momentum today with gains across multiple sectors.",
            ),
            "negative": (
                "Market Concerns Mount",
                "Investors expressed concerns about market conditions today, leading to cautious trading.",
            ),
            "neutral": (
                "Markets Show Mixed Performance",
                "Financial markets showed mixed performance today with no clear directional movement.",
            ),
        }

        return generic_templates.get(
            sentiment,
            ("Market Update", "Financial markets showed varied activity today."),
        )

    def connect(self) -> bool:
        """Simulate connection to mock data source."""
        logger.info("Connecting to mock financial news source")
        self._connected = True
        return True

    def disconnect(self) -> None:
        """Simulate disconnection from mock data source."""
        logger.info("Disconnecting from mock financial news source")
        self._connected = False

    def fetch_news(self, limit: Optional[int] = None) -> List[NewsItem]:
        """
        Fetch mock news items.

        Args:
            limit: Maximum number of news items to fetch

        Returns:
            List[NewsItem]: List of mock news items
        """
        if not self._connected:
            logger.warning("Mock data source not connected")
            return []

        fetch_limit = limit or self.news_count
        news_items = self._mock_news[:fetch_limit]

        logger.info(f"Fetched {len(news_items)} mock news items")
        return news_items

    def is_connected(self) -> bool:
        """Check if mock data source is connected."""
        return self._connected

    @property
    def source_name(self) -> str:
        """Get the name of the mock data source."""
        return f"Mock Financial News (Ollama + {self.model_name})"

    def _generate_mock_data(self) -> None:
        """Generate realistic mock financial news data using LLM or fallback."""
        sources = [
            "Reuters",
            "Bloomberg",
            "Financial Times",
            "Wall Street Journal",
            "CNBC",
        ]
        categories = ["stocks", "bonds", "forex", "commodities", "crypto", "economy"]
        sentiments = ["positive", "negative", "neutral"]

        base_time = datetime.now()

        for i in range(self.news_count):
            # Generate unique ID
            news_id = str(uuid.uuid4())

            # Select random components
            source = random.choice(sources)
            category = random.choice(categories)
            sentiment = random.choice(sentiments)

            # Generate content using LLM or fallback
            if self.use_llm:
                title, content = self._generate_llm_content(category, sentiment, source)
            else:
                title, content = self._generate_fallback_content(category, sentiment)

            # Generate realistic timestamp (within last 24 hours)
            hours_ago = random.randint(0, 24)  # More random timing
            published_date = base_time - timedelta(hours=hours_ago)

            # Create news item
            news_item = NewsItem(
                id=news_id,
                title=title,
                content=content,
                source=source,
                published_date=published_date,
                category=category,
                sentiment=sentiment,
                url=f"https://mock-news.com/article/{news_id}",
                metadata={
                    "read_time": f"{random.randint(2, 8)} min",
                    "author": f"Mock Author {random.randint(1, 100)}",
                    "tags": [category, source.lower(), sentiment],
                    "generated_by": f"ollama_{self.model_name}"
                    if self.use_llm
                    else "fallback_templates",
                },
            )
            logger.info(f"Generated news item: {news_item}")
            self._mock_news.append(news_item)

        logger.info(
            f"Generated {len(self._mock_news)} mock news items using {'Ollama LLM' if self.use_llm else 'fallback templates'}"
        )
