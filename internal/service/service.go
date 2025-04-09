package service

import (
	"context"
	"fmt"

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
