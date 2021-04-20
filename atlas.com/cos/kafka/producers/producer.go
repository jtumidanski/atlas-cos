package producers

import (
	"atlas-cos/rest/requests"
	"atlas-cos/retry"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func createKey(key int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(key))
	return b
}

func produceEvent(l log.FieldLogger, topicToken string, key []byte, event interface{}) {
	td, err := requests.Topic(l).GetTopic(topicToken)
	if err != nil {
		l.WithError(err).Fatalf("Unable to retrieve topic %s for producer.", topicToken)
		return
	}
	name := td.Attributes.Name

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	r, err := json.Marshal(event)
	l.WithField("message", string(r)).Debugf("Writing message to topic %s.", name)
	if err != nil {
		l.WithError(err).Fatalf("Unable to marshall event for topic %s.", name)
	}

	writeMessage := func(attempt int) (bool, error) {
		err = w.WriteMessages(context.Background(), kafka.Message{
			Key:   key,
			Value: r,
		})
		if err != nil {
			l.Warnf("Unable to emit event on topic %s, will retry.", td.Attributes.Name)
			return true, err
		}
		return false, err
	}

	err = retry.Retry(writeMessage, 10)
	if err != nil {
		l.WithError(err).Fatalf("Unable to emit event on topic %s.", td.Attributes.Name)
	}
}
