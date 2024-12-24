package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/s21platform/search-proto/search"
	"github.com/samber/lo"
)

type Handler struct {
	search.UnimplementedSearchServiceServer
	uS userService
}

func New(uS userService) *Handler {
	return &Handler{uS: uS}
}

func (h *Handler) GetSociety(ctx context.Context, in *search.GetSocietyIn) (*search.GetSocietyOut, error) {
	societies := []*search.Society{
		{
			Name:        "Тестовое сообщество 1",
			Description: "Моковое сообщество",
			IsPrivate:   false,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
		{
			Name:        "Тестовое сообщество 2",
			Description: "Моковое сообщество 2",
			IsPrivate:   false,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
		{
			Name:        "Тестовое сообщество 3",
			Description: "Моковое сообщество 3",
			IsPrivate:   true,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
		{
			Name:        "Тестовое сообщество 4",
			Description: "Моковое сообщество 4",
			IsPrivate:   false,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
		{
			Name:        "Тестовое сообщество 5",
			Description: "Моковое сообщество 5",
			IsPrivate:   false,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
		{
			Name:        "Тестовое сообщество 6",
			Description: "Моковое сообщество 6",
			IsPrivate:   false,
			AvatarLink:  "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg",
		},
	}

	if len(in.PartName) > 0 {
		societies = lo.Filter(societies, func(society *search.Society, _ int) bool {
			return strings.Contains(strings.ToLower(society.Name), strings.ToLower(in.PartName))
		})
	}
	total := int64(len(societies))
	start := in.Offset
	end := in.Offset + in.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &search.GetSocietyOut{
		Societies: societies[start:end],
		Total:     total,
	}, nil
}

func (h *Handler) GetUserWithLimit(ctx context.Context, in *search.GetUserWithLimitIn) (*search.GetUserWithLimitOut, error) {
	userOffsetOut, err := h.uS.GetUserWithOffset(ctx, in.Limit, in.Offset, in.Nickname)

	if err != nil {
		return nil, fmt.Errorf("error in GetUserWithOffset: %w", err)
	}

	var usersOut []*search.UserSr
	for _, user := range userOffsetOut.User {
		usersOut = append(usersOut, &search.UserSr{
			Nickname:   user.Nickname,
			Uuid:       user.Uuid,
			AvatarLink: user.AvatarLink,
			Name:       user.Name,
			Surname:    user.Surname,
		})
	}
	return &search.GetUserWithLimitOut{Users: usersOut, Total: userOffsetOut.Total}, nil
}
