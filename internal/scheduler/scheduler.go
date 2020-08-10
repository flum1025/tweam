package scheduler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/flum1025/tweam/internal/aws"
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
	sqsd, err := s.sqsd(workerEndpoint)
	if err != nil {
		return fmt.Errorf("initialize sqsd: %w", err)
	}

	s.cron.AddFunc("@every 1s", sqsd)

	// for _, account := range s.config.Accounts {
	// 	s.cron.AddFunc(fmt.Sprintf("@every %ds", account.FetchInterval), s.directmessage(account))
	// 	s.cron.AddFunc(fmt.Sprintf("@every %ds", account.FetchInterval), s.homeTimeline(account))
	// }

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

func (s *scheduler) sqsd(workerEndpoint string) (func(), error) {
	sqsClient, err := aws.NewSQSClient(s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQSClient: %w", err)
	}

	return func() {
		log.Println("Running sqsd")

		messages, err := sqsClient.ReceiveMessages()
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] failed to receive message: %v", err))
			return
		}

		log.Println(fmt.Sprintf("%d messages received", len(messages)))

		for _, message := range messages {
			go func(message *sqs.Message) {

				client := http.Client{Timeout: timeout}

				_, err := client.Post(workerEndpoint, "application/json", bytes.NewBufferString(*message.Body))
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] request to %v: %v", workerEndpoint, err))
					return
				}

				err = sqsClient.DeleteMessage(*message.ReceiptHandle)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] failed to delete message: %v", err))
					return
				}
			}(message)
		}

		log.Println("Finish sqsd")
	}, nil
}
