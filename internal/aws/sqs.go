package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/flum1025/tweam/internal/config"
)

type SQSClient struct {
	client   *sqs.SQS
	queueURL string
}

func NewSQSClient(
	config *config.Config,
) (*SQSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		EndpointResolver: endpoints.ResolverFunc(
			func(string, string, ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
				return endpoints.ResolvedEndpoint{
					URL: config.QueueEndpoint,
				}, nil
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("aws create session: %w", err)
	}

	client := sqs.New(sess)

	return &SQSClient{
		client:   client,
		queueURL: config.QueueUrl,
	}, nil
}

func (c *SQSClient) SendMessage(message string, groupID string, deduplicationID string) error {
	output, err := c.client.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(message),
		QueueUrl:               aws.String(c.queueURL),
		MessageDeduplicationId: aws.String(deduplicationID),
		MessageGroupId:         aws.String(groupID),
	})

	if err != nil {
		return fmt.Errorf("sqs: send message: %w", err)
	}

	if output.MessageId == nil {
		return fmt.Errorf("sqs: empty message id")
	}

	return nil
}

func (c *SQSClient) ReceiveMessages() ([]*sqs.Message, error) {
	output, err := c.client.ReceiveMessage(
		&sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(10),
			QueueUrl:            aws.String(c.queueURL),
			WaitTimeSeconds:     aws.Int64(20),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("sqs: receive message: %w", err)
	}

	if len(output.Messages) == 0 {
		return nil, nil
	}

	return output.Messages, nil
}

func (c *SQSClient) DeleteMessage(receiptHandle string) error {
	if _, err := c.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}); err != nil {
		return fmt.Errorf("sqs: delete message: %w", err)
	}

	return nil
}
