package kafka

import (
	"news-fetcher/src/fetcher"
	"time"
)

// ConvertToKafkaMessage converts a NewsArticle to a KafkaMessage
func ConvertToKafkaMessage(article fetcher.NewsArticle, fetcherName string) KafkaMessage {
	return KafkaMessage{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		Summary:     article.Summary,
		URL:         article.URL,
		ImageURL:    article.ImageURL,
		Source:      article.Source,
		Author:      article.Author,
		PublishedAt: article.PublishedAt,
		Category:    article.Category,
		Tags:        article.Tags,
		FetchedAt:   time.Now(),
		FetcherName: fetcherName,
	}
}

// ConvertArticlesToKafkaMessages converts multiple NewsArticles to KafkaMessages
func ConvertArticlesToKafkaMessages(articles []fetcher.NewsArticle, fetcherName string) []KafkaMessage {
	messages := make([]KafkaMessage, 0, len(articles))
	for _, article := range articles {
		messages = append(messages, ConvertToKafkaMessage(article, fetcherName))
	}
	return messages
}
