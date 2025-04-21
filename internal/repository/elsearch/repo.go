package elsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

func (c *Client) ExistOrCreateIndex(ctx context.Context, users []string) error {
	_ = ctx
	res, err := c.es.Indices.Exists(users)
	if err != nil {
		return fmt.Errorf("failed to exist index: %s", err)
	}
	if res.StatusCode == http.StatusOK {
		fmt.Println("Index EXIST")
		return nil
	}

	mapping := `{
        "mappings": {
            "properties": {
                "login": { "type": "keyword" },
                "avatar": { "type": "keyword" },
                "name": { "type": "text", "fields": { "raw": { "type": "keyword" } } },
                "surname": { "type": "text", "fields": { "raw": { "type": "keyword" } } },
                "birthdate": { "type": "date", "format": "strict_date_optional_time||epoch_millis" },
                "phone": { "type": "keyword" },
                "telegram": { "type": "keyword" },
                "git": { "type": "keyword" },
                "city": { "type": "long" },
                "os": { "type": "long" },
                "work": { "type": "long" },
                "university": { "type": "long" },
                "uuid": { "type": "keyword" }
            }
        }
    }`

	res, err = c.es.Indices.Create("users", c.es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return fmt.Errorf("failed to create index: %s, err : %s", res.String(), err)
	}

	return nil
}

func (c *Client) BulkIndexUsers(ctx context.Context, users []model.UserInfo) error {
	var buf bytes.Buffer

	for _, u := range users {
		if u.UUID == nil || *u.UUID == "" {
			continue
		}

		meta := fmt.Sprintf(`{ "index" : { "_index" : "users", "_id" : "%s" } }`, *u.UUID)
		buf.WriteString(meta + "\n")

		data, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("marshal user %s failed: %w", *u.UUID, err)
		}
		buf.Write(data)
		buf.WriteByte('\n')
	}

	resp, err := c.es.Bulk(bytes.NewReader(buf.Bytes()), c.es.Bulk.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("bulk index failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("bulk indexing error: %s", resp.String())
	}

	return nil
}
