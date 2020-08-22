package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/flum1025/tweam/internal/app"
	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/entity"
	"github.com/flum1025/tweam/internal/redis"
	"github.com/flum1025/tweam/internal/twitter"
	"github.com/robfig/cron/v3"
)

const (
	timeout                     = time.Duration(25 * time.Second)
	homeTimelineCursorKeyTpl    = "home_timeline_cursor_%s"
	mentionTimelineCursorKeyTpl = "mention_timeline_cursor_%s"
)

type scheduler struct {
	config *config.Config
	cron   *cron.Cron
	app    *app.App
	redis  *redis.RedisClient
}

func NewScheduler(
	config *config.Config,
) (*scheduler, error) {
	app, err := app.NewApp(config)
	if err != nil {
		return nil, fmt.Errorf("get app: %w", err)
	}

	redis, err := redis.NewRedisClient(config)
	if err != nil {
		return nil, fmt.Errorf("get redis client: %w", err)
	}

	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	return &scheduler{
		config: config,
		cron:   c,
		app:    app,
		redis:  redis,
	}, nil
}

func (s *scheduler) Run() {
	for _, account := range s.config.Accounts {
		// FIXME: Cursorだけだと重複排除できない
		// s.cron.AddFunc(fmt.Sprintf("@every %ds", account.DirectmessageFetchInterval), s.directmessage(account))
		s.cron.AddFunc(fmt.Sprintf("@every %ds", account.HomeTimelineFetchInterval), s.homeTimeline(account))
		s.cron.AddFunc(fmt.Sprintf("@every %ds", account.MentionTimeline), s.mentionTimeline(account))
	}

	log.Println(fmt.Sprintf("%d schedules have been registered", len(s.cron.Entries())))

	s.cron.Run()
}

func (s *scheduler) directmessage(config config.Account) func() {
	client := twitter.NewTwitterClient(config.Token)

	return eventWrapper("directmessage", config.ID, func() {
		messages, err := client.DirectMessages("")
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] get direct messages: %v", err))
			return
		}

		if err = s.app.PublishMessages(messages.Messages(config.ID)); err != nil {
			log.Println(fmt.Sprintf("[ERROR] failed to publish messages: %v", err))
			return
		}
	})
}

func (s *scheduler) homeTimeline(config config.Account) func() {
	client := twitter.NewTwitterClient(config.Token)

	return s.timeline(
		"homeTimeline",
		homeTimelineCursorKeyTpl,
		config,
		client.HomeTimeline,
	)
}

func (s *scheduler) mentionTimeline(config config.Account) func() {
	client := twitter.NewTwitterClient(config.Token)

	return s.timeline(
		"mentionTimeline",
		mentionTimelineCursorKeyTpl,
		config,
		client.MentionTimeline,
	)
}

func (s *scheduler) timeline(
	eventTitle string,
	cursorKeyTpl string,
	config config.Account,
	fetcher func(int64) (entity.Tweets, error),
) func() {
	cursorKey := fmt.Sprintf(cursorKeyTpl, config.ID)

	return eventWrapper(eventTitle, config.ID, func() {
		cursorPtr, err := s.redis.GetInt64(cursorKey)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] get cursor key=%s: %v", cursorKey, err))
			return
		}

		var cursor int64
		if cursorPtr != nil {
			cursor = *cursorPtr
		}

		tweets, err := fetcher(cursor)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] get tweet title=%s: %v", eventTitle, err))
			return
		}

		if len(tweets) == 0 {
			return
		}

		if err = s.app.PublishMessages(tweets.Messages(config.ID)); err != nil {
			log.Println(fmt.Sprintf("[ERROR] failed to publish messages: %v", err))
			return
		}

		if err = s.redis.Set(cursorKey, tweets[0].ID, 0); err != nil {
			log.Println(fmt.Sprintf("[ERROR] set cursor key=%s: %v", cursorKey, err))
			return
		}
	})
}

func eventWrapper(title string, userID string, callback func()) func() {
	return func() {
		log.Println(fmt.Sprintf("Running %s, user_id=%v", title, userID))
		callback()
		log.Println(fmt.Sprintf("Finish %s, user_id=%v", title, userID))
	}
}
