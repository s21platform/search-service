package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	search "github.com/s21platform/search-proto/search/new_uuid"
	"github.com/s21platform/search-service/internal/config"
)

type Handler struct {
	els config.Elastic
	uC  UserClient
}

func New(els config.Elastic, uC UserClient) *Handler {
	return &Handler{els: els, uC: uC}
}

func convertMessage(bMessage []byte, target interface{}) error {
	err := json.Unmarshal(bMessage, target)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Handler(ctx context.Context, in []byte) {
	fmt.Println("Message:", string(in), ctx)
	var msg search.UpdateUser
	err := convertMessage(in, &msg)
	if err != nil {
		log.Println("convert err:", err)
		return
	}
	ctx = context.WithValue(ctx, config.KeyUUID, msg.Uuid)
	res, err := h.uC.GetUserInfoByUUID(ctx, msg.Uuid)
	if err != nil {
		log.Println("failed to GetInfo from user-service:", err)
		return
	}
	fmt.Println(res)
}
