package elastic

import (
	"fmt"

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
