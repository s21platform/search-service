package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/s21platform/search-service/internal/model"
	user_proto "github.com/s21platform/user-proto/user-proto"

	search "github.com/s21platform/search-proto/search/new_uuid"
	"github.com/s21platform/search-service/internal/config"
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

func toUserDoc(user *user_proto.GetUserInfoByUUIDOut) model.UserInfo {
	os := model.GetOs{
		Id:    user.Os.Id,
		Label: user.Os.Label,
	}
	return model.UserInfo{
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		Name:       safeString(user.Name),
		Surname:    safeString(user.Surname),
		Birthdate:  safeString(user.Birthdate),
		Phone:      safeString(user.Phone),
		City:       safeString(user.City),
		Telegram:   safeString(user.Telegram),
		Git:        safeString(user.Git),
		Os:         os,
		Work:       safeString(user.Work),
		University: safeString(user.University),
		Skills:     user.Skills,
		Hobbies:    user.Hobbies,
	}
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
	fmt.Println("Get UserInfo:\n", res)

	resForUpdate := toUserDoc(res)
	err = h.els.Update(ctx, msg.Uuid, resForUpdate)
	fmt.Println("err elastic", err)
}
