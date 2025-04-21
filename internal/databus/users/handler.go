package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

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
	return model.UserInfo{
		Nickname:       user.Nickname,
		LastAvatarLink: user.Avatar,
		Name:           user.Name,
		Surname:        user.Surname,
		Birthdate:      parseDate(user.Birthdate),
		Phone:          user.Phone,
		Telegram:       user.Telegram,
		Git:            user.Git,
		CityId:         parseInt64FromString(user.City),
		OSId:           &user.Os.Id,
		WorkId:         parseInt64FromString(user.Work),
		UniversityId:   parseInt64FromString(user.University),
		UUID:           user.Uuid,
	}
}

func parseDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func parseInt64FromString(s *string) *int64 {
	if s == nil {
		return nil
	}
	val, err := strconv.ParseInt(*s, 10, 64)
	if err != nil {
		return nil
	}
	return &val
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
