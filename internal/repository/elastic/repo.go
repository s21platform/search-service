package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/s21platform/search-service/internal/model"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/s21platform/search-service/internal/config"
)

type Client struct {
	es *elasticsearch.Client
}

func MustConnect(cfg *config.Config) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.Elastic.Server},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Println("failed to create client:", err)
		return nil, err
	}

	return &Client{es: es}, nil
}

func (c *Client) SetUser(ctx context.Context, data map[string]interface{}) error {
	uuid, exist := data["Uuid"]
	if !exist {
		return fmt.Errorf("failed to user uuid is empty")
	}
	now := time.Now()
	data["updatedAt"] = now

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	req := esapi.IndexRequest{
		Index:      "users",
		DocumentID: uuid.(string),
		Body:       bytes.NewReader(jsonData),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to execute elasticsearch request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to answer elasticsearch [%s] %s", res.Status(), res.String())
	}

	return nil
}

func (c *Client) SearchUser(ctx context.Context,
	query string,
	filters map[string]interface{},
	page, pageSize int) ([]*model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 20 {
		pageSize = 20
	}
	from := (page - 1) * pageSize

	searchBody, err := c.buildSearchQuery(query, filters, from, pageSize)

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex("user_create"),
		c.es.Search.WithBody(searchBody),
		c.es.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, 0, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search request failed [%s]", res.Status())
	}

	users, total, err := c.parseSearchResponse(res)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse search response %w", err)
	}

	return users, total, nil
}

func (c *Client) buildSearchQuery(query string, filters map[string]interface{}, from, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   c.buildMustClause(query),
				"filter": c.buildFilterClause(filters),
			},
		},
		"from": from,
		"size": size,
		"sort": c.buildSortClause(),
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	return &buf, nil
}

func (c *Client) buildMustClause(query string) map[string]interface{} {
	if query == "" {
		return map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": query,
			"fields": []string{
				"nickname^3",
				"email^2",
				"first_name",
				"last_name",
				"city",
			},
			"fuzziness": "AUTO",
		},
	}
}

func (c *Client) buildFilterClause(filters map[string]interface{}) []map[string]interface{} {
	if filters == nil {
		return nil
	}
	var filterClauses []map[string]interface{}

	for field, value := range filters {
		switch v := value.(type) {
		case string, bool, int, float64:
			filterClauses = append(filterClauses, map[string]interface{}{
				"term": map[string]interface{}{field: v},
			})
		case []string:
			filterClauses = append(filterClauses, map[string]interface{}{
				"terms": map[string]interface{}{field: v},
			})
		case map[string]interface{}:
			if rangeFilter, ok := v["range"]; ok {
				filterClauses = append(filterClauses, map[string]interface{}{
					"range": map[string]interface{}{field: rangeFilter},
				})
			}
		default:
			continue
		}
	}

	return filterClauses
}

func (c *Client) buildSortClause() []map[string]interface{} {
	return []map[string]interface{}{
		{"_score": "desc"},
	}
}

func (c *Client) parseSearchResponse(res *esapi.Response) ([]*model.User, int, error) {
	var searchResult struct {
		Took int `json:"took"` //time to response
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				User  model.User `json:"_source"`
				Relev float64    `json:"_score"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}

	users := make([]*model.User, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		users[i] = &hit.User
	}

	return users, searchResult.Hits.Total.Value, nil
}
