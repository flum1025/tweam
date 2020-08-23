package entity

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

type FollowType string

var (
	FollowTypeFollow   FollowType = "follow"
	FollowTypeUnfollow FollowType = "unfollow"
)

type FollowUser struct {
	twitter.User
	ID                        string `json:"id"`
	ProfileLinkColor          int    `json:"profile_link_color"`
	ProfileTextColor          int    `json:"profile_text_color"`
	ProfileSidebarBorderColor int    `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor   int    `json:"profile_sidebar_fill_color"`
	ProfileBackgroundColor    int    `json:"profile_background_color"`
}

func (u FollowUser) TwitterUser() twitter.User {
	user := u.User

	id, _ := strconv.ParseInt(u.ID, 10, 64)
	user.ID = id
	user.IDStr = u.ID

	return user
}

type FollowRaw struct {
	Type      FollowType `json:"type"`
	CreatedAt string     `json:"created_timestamp"`
	Target    FollowUser `json:"target"`
	Source    FollowUser `json:"source"`
}

type Follow struct {
	Type      FollowType   `json:"type"`
	CreatedAt string       `json:"created_timestamp"`
	Target    twitter.User `json:"target"`
	Source    twitter.User `json:"source"`
}

func (f FollowRaw) NewMessage(userID string) Message {
	return Message{
		ForUserID: userID,
		Type:      MessageTypeFollow,
		Data: Follow{
			Type:      f.Type,
			CreatedAt: f.CreatedAt,
			Target:    f.Target.TwitterUser(),
			Source:    f.Source.TwitterUser(),
		},
	}
}
