package users

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/s21platform/search-service/internal/config"
)

type Handler struct {
	els config.Elastic
}

func New(els config.Elastic) *Handler {
	return &Handler{els: els}
}

func convertMessage(bMessage []byte, target interface{}) error {
	err := json.Unmarshal(bMessage, target)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Handler(ctx context.Context, in []byte) {
	_ = ctx
	var str string
	err := convertMessage(in, &str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("work") ///////////////////////////////////////////////////////del, but not now
	fmt.Println("Message:", str)
}
