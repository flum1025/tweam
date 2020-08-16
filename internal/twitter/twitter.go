package twitter

import (
	"fmt"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/entity"
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

func (c *TwitterClient) DirectMessages(cursor string) (entity.DirectMessages, error) {
	events, _, err := c.client.DirectMessages.EventsList(&twitter.DirectMessageEventsListParams{
		Count:  50,
		Cursor: cursor,
	})

	if err != nil {
		return nil, fmt.Errorf("twitter: get directmessage event list: %w", err)
	}

	userIDs := make([]int64, 0)

	for _, event := range events.Events {
		recipientID, _ := strconv.ParseInt(event.Message.Target.RecipientID, 10, 64)
		senderID, _ := strconv.ParseInt(event.Message.SenderID, 10, 64)

		userIDs = append(userIDs, recipientID)
		userIDs = append(userIDs, senderID)
	}

	users, err := c.Users(userIDs)
	if err != nil {
		return nil, fmt.Errorf("twitter: find users: %w", err)
	}

	userMap := make(map[int64]twitter.User)

	for _, user := range users {
		userMap[user.ID] = user
	}

	ms := make([]entity.DirectMessage, 0, len(events.Events))

	for _, event := range events.Events {
		recipientID, _ := strconv.ParseInt(event.Message.Target.RecipientID, 10, 64)
		senderID, _ := strconv.ParseInt(event.Message.SenderID, 10, 64)

		recipient, ok := userMap[recipientID]
		if !ok {
			return nil, fmt.Errorf("twitter: user not found: %d", recipientID)
		}

		sender, ok := userMap[senderID]
		if !ok {
			return nil, fmt.Errorf("twitter: user not found: %d", senderID)
		}

		ms = append(ms, entity.NewDirectMessageFromDirectMessageEvent(event, &recipient, &sender))
	}

	return ms, nil
}

func (c *TwitterClient) Users(ids []int64) ([]twitter.User, error) {
	users, _, err := c.client.Users.Lookup(&twitter.UserLookupParams{
		UserID:          ids,
		IncludeEntities: twitter.Bool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("twtiter: get users: %w", err)
	}

	return users, nil
}

func (c *TwitterClient) HomeTimeline(cursor int64) (entity.Tweets, error) {
	tweets, _, err := c.client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count:          200,
		ExcludeReplies: twitter.Bool(false),
		SinceID:        cursor,
	})

	if err != nil {
		return nil, fmt.Errorf("twitter: get home_timeline: %w", err)
	}

	ms := make([]entity.Tweet, 0, len(tweets))

	for _, tweet := range tweets {
		ms = append(ms, entity.Tweet(tweet))
	}

	return ms, nil
}

func (c *TwitterClient) MentionTimeline(cursor int64) (entity.Tweets, error) {
	tweets, _, err := c.client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
		Count:   200,
		SinceID: cursor,
	})

	if err != nil {
		return nil, fmt.Errorf("twitter: get mention_timeline %w", err)
	}

	ms := make([]entity.Tweet, 0, len(tweets))

	for _, tweet := range tweets {
		ms = append(ms, entity.Tweet(tweet))
	}

	return ms, nil
}
