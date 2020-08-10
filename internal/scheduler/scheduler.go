package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/twitter"
	"github.com/robfig/cron/v3"
)

const timeout = time.Duration(25 * time.Second)

type scheduler struct {
	config *config.Config
	cron   *cron.Cron
}

func NewScheduler(
	config *config.Config,
) *scheduler {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	return &scheduler{
		config: config,
		cron:   c,
	}
}

func (s *scheduler) Run(workerEndpoint string) error {
	for _, account := range s.config.Accounts {
		s.cron.AddFunc(fmt.Sprintf("@every %ds", account.DirectmessageFetchInterval), s.directmessage(account))
		s.cron.AddFunc(fmt.Sprintf("@every %ds", account.HomeTimelineFetchInterval), s.homeTimeline(account))
	}

	log.Println(fmt.Sprintf("%d schedules have been registered", len(s.cron.Entries())))

	s.cron.Run()

	return nil
}

func (s *scheduler) directmessage(config config.Account) func() {
	_ = twitter.NewTwitterClient(config.Token)

	return func() {
		log.Println("Running directmessage")
		log.Println("Finish directmessage")
	}
}

func (s *scheduler) homeTimeline(config config.Account) func() {
	return func() {
		log.Println("Running homeTimeline")
		log.Println("Finish homeTimeline")
	}
}
