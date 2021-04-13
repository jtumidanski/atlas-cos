package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type gainExperienceEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

func GainExperienceEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &gainExperienceEvent{}
	}
}

func HandleGainExperienceEvent(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*gainExperienceEvent); ok {
			character.Processor(l, db).GainExperience(event.CharacterId, event.PersonalGain + event.PartyGain)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [GainExperienceEvent]")
		}
	}
}
