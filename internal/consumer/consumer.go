package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/nogavadu/notification-service/internal/clients/kafka"
	"github.com/nogavadu/notification-service/internal/service"
)

type consumer struct {
	ctx    context.Context
	cancel context.CancelFunc

	brokers []string
	topics  []string
	groupId string

	client sarama.ConsumerGroup

	emailServ service.EmailService
}

func New(
	brokers []string,
	topics []string,
	groupId string,
	emailServ service.EmailService,
) (kafka.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.AutoCommit.Enable = false

	client, err := sarama.NewConsumerGroup(brokers, groupId, cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &consumer{
		ctx:       ctx,
		cancel:    cancel,
		brokers:   brokers,
		topics:    topics,
		groupId:   groupId,
		client:    client,
		emailServ: emailServ,
	}, nil
}

func (c *consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			switch msg.Topic {
			case "registrations-topic":
				c.emailServ.SendMsg(
					session.Context(),
					[]string{string(msg.Value)},
					"Спасибо за регистрацию",
					"Успешная регистрация",
				)
			}

			session.MarkMessage(msg, "")
			session.Commit()
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *consumer) Start() error {
	for {
		if err := c.client.Consume(c.ctx, c.topics, c); err != nil {
			return err
		}

		if c.ctx.Err() != nil {
			return nil
		}
	}
}

func (c *consumer) Stop() error {
	c.cancel()
	return c.client.Close()
}
