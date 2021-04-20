package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func CharacterStatusEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterStatusEvent{}
	}
}

func HandleCharacterStatusEvent(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			l.Debugf("Begin event handling.")
			if event.Type == "LOGIN" {
				character.Processor(l, db).UpdateLoginPosition(event.CharacterId)
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
