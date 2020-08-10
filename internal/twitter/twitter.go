package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/flum1025/tweam/internal/config"
)

type TwitterClient struct {
	client *twitter.Client
}

func NewTwitterClient(
	token config.Token,
) *TwitterClient {
	c := oauth1.NewConfig(token.ConsumerKey, token.ConsumerSecret)
	t := oauth1.NewToken(token.AccessToken, token.AccessTokenSecret)
	httpClient := c.Client(oauth1.NoContext, t)
	client := twitter.NewClient(httpClient)

	return &TwitterClient{
		client: client,
	}
}
