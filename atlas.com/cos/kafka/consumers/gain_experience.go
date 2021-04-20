package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*gainExperienceEvent); ok {
			l.Debugf("Begin event handling.")
			character.Processor(l, db).GainExperience(event.CharacterId, event.PersonalGain + event.PartyGain)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
