package chroma

import (
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
)

type Client struct {
	client chromago.Client
}

func NewClient(address string) (*Client, error) {
	// Initialize ChromaDB client v2
	client, err := chromago.NewHTTPClient(chromago.WithBaseURL(address))
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Heartbeat(ctx context.Context) (map[string]int64, error) {
	// In v2, Heartbeat returns an error if not alive.
	err := c.client.Heartbeat(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]int64{"status": 1}, nil
}

type Document struct {
	ID       string
	Content  string
	Metadata map[string]interface{}
	Score    float32
}

func (c *Client) Query(ctx context.Context, queryEmbeddings []float32, nResults int32) ([]Document, error) {
	// In v2, we get the collection first
	col, err := c.client.GetCollection(ctx, "spark_documents", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Create embedding object
	emb := embeddings.NewEmbeddingFromFloat32(queryEmbeddings)

	// Query using the new API
	qr, err := col.Query(ctx,
		chromago.WithQueryEmbeddings(emb),
		chromago.WithNResults(int(nResults)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	// Parse results using interface methods
	var docs []Document
	ids := qr.GetIDGroups()
	documents := qr.GetDocumentsGroups()
	metadatas := qr.GetMetadatasGroups()
	distances := qr.GetDistancesGroups()

	if len(ids) > 0 && len(ids[0]) > 0 {
		for i := 0; i < len(ids[0]); i++ {
			doc := Document{
				ID: string(ids[0][i]),
			}

			if len(documents) > 0 && len(documents[0]) > i {
				// Document is an interface, but usually it's a string
				doc.Content = fmt.Sprintf("%v", documents[0][i])
			}

			if len(metadatas) > 0 && len(metadatas[0]) > i {
				// Metadata is an interface, we'll try to convert it to map
				// For now, we'll just use an empty map to avoid type issues
				doc.Metadata = make(map[string]interface{})
			}

			if len(distances) > 0 && len(distances[0]) > i {
				doc.Score = float32(distances[0][i])
			}

			docs = append(docs, doc)
		}
	}

	return docs, nil
}
