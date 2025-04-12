package elsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/model"
)

type Client struct {
	es *elasticsearch.Client
}

func New(cfg config.Elastic) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.Server},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Println("failed to create client:", err)
		return nil, err
	}

	return &Client{es: es}, nil
}

func (c *Client) Update(ctx context.Context, id string, doc model.UserInfo) error {
	body, err := json.Marshal(map[string]interface{}{
		"doc":           doc,
		"doc_as_upsert": true,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	resp, err := c.es.Update(
		"users",
		id,
		bytes.NewReader(body),
		c.es.Update.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("error updating document: %s", resp.String())
	}

	return nil
}
