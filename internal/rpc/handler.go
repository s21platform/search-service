package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/s21platform/search-proto/search"
	users "github.com/s21platform/user-proto/user-proto"
	"github.com/samber/lo"
)

type Handler struct {
	search search.UnimplementedSearchServiceServer
	users  users.UnimplementedUserServiceServer
}

func New() *Handler {
	return &Handler{}
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
	usersUs, err := h.users.GetUserWithOffset(ctx, &users.GetUserWithOffsetIn{
		Limit:    in.Limit,
		Offset:   in.Offset,
		Nickname: in.Nickname,
	})
	if err != nil {
		return nil, fmt.Errorf("error h.users.GetUserWithOffset: %w", err)
	}
	var usersSr []*search.User
	for _, u := range usersUs.User {
		usersSr = append(usersSr, &search.User{
			Nickname:   u.Nickname,
			Uuid:       u.Uuid,
			AvatarLink: u.AvatarLink,
			Name:       u.Name,
			Surname:    u.Surname,
		})
	}
	return &search.GetUserWithLimitOut{
		Users: usersSr,
		Total: usersUs.Total,
	}, nil
}
