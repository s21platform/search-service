package user_update

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/model"
	"github.com/s21platform/user-service/pkg/user"
	"log"
)

type Handler struct {
	els Elastic
	uC  UserClient
}

func New(els Elastic, uC UserClient) *Handler {
	return &Handler{els: els, uC: uC}
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
		log.Printf("failed to convert message: %v", err)
		return err
	}

	user := model.User{
		UUID:      msg.UserUuid,
		Nickname:  msg.UserNickname,
		Email:     msg.UserEmail,
		FirstName: "Aleksey", //заглушка
		LastName:  "Popov",   //заглушка
	}

	err = h.els.SetUser(ctx, &user)
	if err != nil {
		m.Increment("user_create.error")
		return fmt.Errorf("failed to set user: %v", err)
	}
	return nil
}
