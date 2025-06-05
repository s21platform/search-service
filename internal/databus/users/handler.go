package users

import (
	"context"
	"fmt"
)

type Handler struct {
	els Elastic
	uC  UserClient
}

func New(els Elastic, uC UserClient) *Handler {
	return &Handler{els: els, uC: uC}
}

func (h *Handler) Handler(ctx context.Context, in []byte) error {
	fmt.Println("Message:", string(in), ctx)
	return nil
}
