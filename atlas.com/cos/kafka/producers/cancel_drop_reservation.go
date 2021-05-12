package producers

import (
	"github.com/sirupsen/logrus"
)

type cancelDropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func CancelDropReservation(l logrus.FieldLogger) func(dropId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND")
	return func(dropId uint32, characterId uint32) {
		event := &cancelDropReservationEvent{characterId, dropId}
		producer(CreateKey(int(dropId)), event)
	}
}