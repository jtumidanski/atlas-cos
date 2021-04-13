package consumers

import (
	"atlas-cos/drop"
	"gorm.io/gorm"
	"log"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func DropReservationEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &dropReservationEvent{}
	}
}

func HandleDropReservationEvent(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*dropReservationEvent); ok {
			if event.Type == "SUCCESS" {
				drop.Processor(l, db).AttemptPickup(event.CharacterId, event.DropId)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [DropReservationEvent]")
		}
	}
}
