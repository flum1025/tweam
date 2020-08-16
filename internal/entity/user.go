package entity

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

type EventUser struct {
	ID                   string `json:"id"`
	CreatedTimestamp     string `json:"created_timestamp"`
	Name                 string `json:"name"`
	ScreenName           string `json:"screen_name"`
	Description          string `json:"description"`
	Protected            bool   `json:"protected"`
	Verified             bool   `json:"verified"`
	FollowersCount       int    `json:"followers_count"`
	FriendsCount         int    `json:"friends_count"`
	StatusesCount        int    `json:"statuses_count"`
	ProfileImageURL      string `json:"profile_image_url"`
	ProfileImageURLHTTPS string `json:"profile_image_url_https"`
}

func (u EventUser) NewUser() *twitter.User {
	id, _ := strconv.ParseInt(u.ID, 10, 64)

	return &twitter.User{
		ID:                   id,
		CreatedAt:            u.CreatedTimestamp,
		Name:                 u.Name,
		ScreenName:           u.ScreenName,
		Description:          u.Description,
		Protected:            u.Protected,
		Verified:             u.Verified,
		FollowersCount:       u.FollowersCount,
		FriendsCount:         u.FriendsCount,
		StatusesCount:        u.StatusesCount,
		ProfileImageURL:      u.ProfileImageURL,
		ProfileImageURLHttps: u.ProfileImageURLHTTPS,
	}
}
