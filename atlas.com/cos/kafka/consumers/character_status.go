package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
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
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				character.Processor(l, db).UpdateLoginPosition(event.CharacterId)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [CharacterStatusEvent]")
		}
	}
}
