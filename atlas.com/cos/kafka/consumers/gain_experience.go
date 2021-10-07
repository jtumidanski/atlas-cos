package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
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

func GainExperienceEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &gainExperienceEvent{}
	}
}

func HandleGainExperienceEvent(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*gainExperienceEvent); ok {
			l.Debugf("Begin event handling.")
			character.GainExperience(l, db, span)(event.CharacterId, event.PersonalGain + event.PartyGain)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
