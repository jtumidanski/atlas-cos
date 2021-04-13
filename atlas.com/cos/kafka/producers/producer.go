package producers

import (
	"atlas-cos/rest/requests"
	"atlas-cos/retry"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

func createKey(key int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(key))
	return b
}

func produceEvent(l *log.Logger, topicToken string, key []byte, event interface{}) {
	td, err := requests.Topic(l).GetTopic(topicToken)
	if err != nil {
		l.Fatal("[ERROR] unable to retrieve topic %s for producer.", topicToken)
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        td.Attributes.Name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	r, err := json.Marshal(event)
	if err != nil {
		l.Fatal("[ERROR] unable to marshall event for topic %s with reason %s", td.Attributes.Name, err.Error())
	}

	writeMessage := func(attempt int) (bool, error) {
		err = w.WriteMessages(context.Background(), kafka.Message{
			Key:   key,
			Value: r,
		})
		if err != nil {
			l.Printf("[WARN] unable to emit event on topic %s, will retry.", td.Attributes.Name)
			return true, err
		}
		return false, err
	}

	err = retry.Retry(writeMessage, 10)
	if err != nil {
		l.Fatalf("[ERROR] unable to emit event on topic %s, with reason %s", td.Attributes.Name, err.Error())
	}
}


