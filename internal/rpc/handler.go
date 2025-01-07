package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/search-proto/search"

	"github.com/s21platform/search-service/internal/config"
)

type Handler struct {
	search.UnimplementedSearchServiceServer
	uS userService
	fS friendsService
	sS societyService
}

func New(uS userService, fS friendsService, sS societyService) *Handler {
	return &Handler{uS: uS, fS: fS, sS: sS}
}

func (h *Handler) GetSociety(ctx context.Context, in *search.GetSocietyIn) (*search.GetSocietyOut, error) {
	societies := []*search.SearchSociety{
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
		societies = lo.Filter(societies, func(society *search.SearchSociety, _ int) bool {
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
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetUserWithLimit")
	userOffsetOut, err := h.uS.GetUserWithOffset(ctx, in.Limit, in.Offset, in.Nickname)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user wiht offset: %v", err))
		return nil, fmt.Errorf("error in GetUserWithOffset: %w", err)
	}

	var usersOut []*search.UserSr
	for _, user := range userOffsetOut.User {
		isFriend, err := h.fS.IsFriendsExist(ctx, user.Uuid)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get user friend: %v", err))
		}
		usersOut = append(usersOut, &search.UserSr{
			Nickname:   user.Nickname,
			Uuid:       user.Uuid,
			AvatarLink: user.AvatarLink,
			Name:       user.Name,
			Surname:    user.Surname,
			IsFriend:   isFriend,
		})
	}
	return &search.GetUserWithLimitOut{Users: usersOut, Total: userOffsetOut.Total}, nil
}

func (h *Handler) GetSocietyWithLimit(ctx context.Context, in *search.GetSocietyWithLimitIn) (*search.GetSocietyWithLimitOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSocietyWithLimit")
	societyOffsetOut, err := h.sS.GetSocietyWithOffset(ctx, in.Limit, in.Offset, in.Name)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get society wiht offset: %v", err))
		return nil, fmt.Errorf("error in GetSocietyWithOffset: %w", err)
	}
	var societiesOut []*search.SocietySr
	for _, society := range societyOffsetOut.Society {
		societiesOut = append(societiesOut, &search.SocietySr{
			Name:       society.Name,
			AvatarLink: society.AvatarLink,
			SocietyId:  society.SocietyId,
			IsMember:   society.IsMember,
			IsPrivate:  society.IsPrivate,
		})
	}
	return &search.GetSocietyWithLimitOut{Societies: societiesOut, Total: societyOffsetOut.Total}, nil
}
