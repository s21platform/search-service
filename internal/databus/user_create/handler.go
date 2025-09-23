package user_create

import (
	"context"
	"encoding/json"
	"fmt"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/user-service/pkg/user"
	"log"
)

type Handler struct {
	els Elastic
}

func New(els Elastic) *Handler {
	return &Handler{els: els}
}

func convertMessage(bMessage []byte, target interface{}) error {
	err := json.Unmarshal(bMessage, target)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Handler(ctx context.Context, in []byte) error {
	fmt.Println("Message:", string(in), ctx)

	m := pkg.FromContext(ctx, config.KeyMetrics)

	var msg user.UserCreatedMessage
	err := convertMessage(in, &msg)
	if err != nil {
		m.Increment("user_create.error")
		logger_lib.Error(ctx, fmt.Sprintf("failed to convert message: %v", err))
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(msg.RawMessage, &data)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	err = h.els.SetUser(ctx, data)
	if err != nil {
		m.Increment("user_create.error")
		logger_lib.Error(ctx, fmt.Sprintf("failed to set user in elastic: %v", err))
		return fmt.Errorf("failed to set user in elastic: %v", err)
	}
	return nil
}
