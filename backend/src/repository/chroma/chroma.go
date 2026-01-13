package chroma

import (
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

type Client struct {
	client *chromago.Client
}

func NewClient(address string) (*Client, error) {
	// Initialize ChromaDB client
	// Using NewClient from the library which accepts the address directly or via options depending on version.
	// Based on the error "undefined: chromago.WithBasePath", it seems the version we have uses a different way.
	// Let's try passing the address directly if supported, or check the library source if possible.
	// For v0.1.2, it might be NewClient(basePath string)
	client, err := chromago.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Heartbeat(ctx context.Context) (map[string]int64, error) {
	// The error said: cannot use c.client.Heartbeat(ctx) (value of type map[string]float32) as type map[string]int64
	hb, err := c.client.Heartbeat(ctx)
	if err != nil {
		return nil, err
	}

	// Convert map[string]float32 to map[string]int64
	res := make(map[string]int64)
	for k, v := range hb {
		res[k] = int64(v)
	}
	return res, nil
}

type Document struct {
	ID       string
	Content  string
	Metadata map[string]interface{}
	Score    float32
}

func (c *Client) Query(ctx context.Context, queryEmbeddings []float32, nResults int32) ([]Document, error) {
	col, err := c.client.GetCollection(ctx, "spark_documents", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Create embedding object
	emb := types.NewEmbeddingFromFloat32(queryEmbeddings)

	// Query using QueryWithOptions
	qr, err := col.QueryWithOptions(ctx,
		types.WithQueryEmbeddings([]*types.Embedding{emb}),
		types.WithNResults(nResults),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	// Parse results
	var docs []Document
	if len(qr.Ids) > 0 && len(qr.Ids[0]) > 0 {
		count := len(qr.Ids[0])
		for i := 0; i < count; i++ {
			doc := Document{
				ID: qr.Ids[0][i],
			}

			if len(qr.Documents) > 0 && len(qr.Documents[0]) > i {
				doc.Content = qr.Documents[0][i]
			}

			if len(qr.Metadatas) > 0 && len(qr.Metadatas[0]) > i {
				doc.Metadata = qr.Metadatas[0][i]
			}

			if len(qr.Distances) > 0 && len(qr.Distances[0]) > i {
				doc.Score = qr.Distances[0][i]
			}

			docs = append(docs, doc)
		}
	}

	return docs, nil
}
