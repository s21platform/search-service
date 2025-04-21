package model

import (
	"strconv"
	"time"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

type UserInfo struct {
	Nickname       string     `json:"login"`
	LastAvatarLink string     `json:"avatar"`
	Name           *string    `json:"name"`
	Surname        *string    `json:"surname"`
	Birthdate      *time.Time `json:"birthdate"`
	Phone          *string    `json:"phone"`
	Telegram       *string    `json:"telegram"`
	Git            *string    `json:"git"`
	CityId         *int64     `json:"city"`
	OSId           *int64     `json:"os"`
	WorkId         *int64     `json:"work"`
	UniversityId   *int64     `json:"university"`
	UUID           *string    `json:"uuid"`
}

func ToUserInfoList(in *user_proto.GetUserWithOffsetOutAll) []UserInfo {
	var result []UserInfo

	for _, user := range in.User {
		var ui UserInfo

		ui.Nickname = user.GetNickname()
		ui.LastAvatarLink = user.GetAvatar()

		if user.Name != nil {
			ui.Name = user.Name
		}
		if user.Surname != nil {
			ui.Surname = user.Surname
		}
		if user.Birthdate != nil {
			if birthdate, err := time.Parse("2006-01-02", user.GetBirthdate()); err == nil {
				ui.Birthdate = &birthdate
			}
		}
		if user.Phone != nil {
			ui.Phone = user.Phone
		}
		if user.Telegram != nil {
			ui.Telegram = user.Telegram
		}
		if user.Git != nil {
			ui.Git = user.Git
		}
		if user.City != nil {
			if cityId, err := strconv.ParseInt(user.GetCity(), 10, 64); err == nil {
				ui.CityId = &cityId
			}
		}
		if user.Os != nil {
			ui.OSId = &user.Os.Id
		}
		if user.Work != nil {
			if workId, err := strconv.ParseInt(user.GetWork(), 10, 64); err == nil {
				ui.WorkId = &workId
			}
		}
		if user.University != nil {
			if universityId, err := strconv.ParseInt(user.GetUniversity(), 10, 64); err == nil {
				ui.UniversityId = &universityId
			}
		}
		if user.Uuid != nil {
			ui.UUID = user.Uuid
		}

		result = append(result, ui)
	}

	return result
}

type GetOs struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}
