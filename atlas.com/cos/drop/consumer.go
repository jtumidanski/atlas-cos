package drop

import (
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameReservation = "drop_reservation_event"
	topicTokenReservation   = "TOPIC_DROP_RESERVATION_EVENT"
)

func ReservationConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[reservationEvent](consumerNameReservation, topicTokenReservation, groupId, handleReservationEvent(db))
	}
}

type reservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func handleReservationEvent(db *gorm.DB) kafka.HandlerFunc[reservationEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event reservationEvent) {
		if event.Type == "SUCCESS" {
			AttemptPickup(l, db, span)(event.CharacterId, event.DropId)
		}
	}
}
