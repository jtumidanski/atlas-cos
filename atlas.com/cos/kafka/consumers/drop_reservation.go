package consumers

import (
	"atlas-cos/drop"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func DropReservationEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &dropReservationEvent{}
	}
}

func HandleDropReservationEvent(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*dropReservationEvent); ok {
			l.Debugf("Begin event handling.")
			if event.Type == "SUCCESS" {
				drop.AttemptPickup(l, db)(event.CharacterId, event.DropId)
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
