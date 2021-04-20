package consumers

import (
	"atlas-cos/rest/requests"
	"atlas-cos/retry"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Consumer struct {
	l                 *log.Logger
	ctx               context.Context
	groupId           string
	topicToken        string
	emptyEventCreator EmptyEventCreator
	h                 EventProcessor
}

func NewConsumer(l *log.Logger, ctx context.Context, h EventProcessor, options ...ConsumerOption) Consumer {
	c := &Consumer{}
	c.l = l
	c.ctx = ctx
	c.h = h
	for _, option := range options {
		option(c)
	}
	return *c
}

type EmptyEventCreator func() interface{}

type EventProcessor func(*log.Logger, interface{})

type ConsumerOption func(c *Consumer)

func SetGroupId(groupId string) func(c *Consumer) {
	return func(c *Consumer) {
		c.groupId = groupId
	}
}

func SetTopicToken(topicToken string) func(c *Consumer) {
	return func(c *Consumer) {
		c.topicToken = topicToken
	}
}

func SetEmptyEventCreator(f EmptyEventCreator) func(c *Consumer) {
	return func(c *Consumer) {
		c.emptyEventCreator = f
	}
}

func (c Consumer) Init() {
	td, err := requests.Topic(c.l).GetTopic(c.topicToken)
	if err != nil {
		c.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
		return
	}

	c.l.Printf("[INFO] creating topic consumer for %s", td.Attributes.Name)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: c.groupId,
		MaxWait: 50 * time.Millisecond,
	})

	for {
		msg, err := retry.RetryResponse(consumerReader(c.l, r, c.ctx), 10)
		if err != nil {
			c.l.Fatalf("[ERROR] could not successfully read message " + err.Error())
		}
		if val, ok := msg.(*kafka.Message); ok {
			event := c.emptyEventCreator()
			err = json.Unmarshal(val.Value, &event)
			if err != nil {
				c.l.Println("[ERROR] could not unmarshal event into event class ", val.Value)
			} else {
				c.h(c.l, event)
			}
		}
	}
}

func consumerReader(l *log.Logger, r *kafka.Reader, ctx context.Context) retry.RetryResponseFunc {
	return func(attempt int) (bool, interface{}, error) {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			l.Printf("[WARN] could not successfully read message on topic %s, will retry", r.Config().Topic)
			return true, nil, err
		}
		return false, &msg, err
	}
}
