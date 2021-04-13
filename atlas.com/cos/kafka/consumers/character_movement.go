package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type characterMovementEvent struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement rawMovement `json:"rawMovement"`
}

type rawMovement []byte

func CharacterMovementEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterMovementEvent{}
	}
}

func HandleCharacterMovementEvent(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*characterMovementEvent); ok {
			if event.X != 0 || event.Y != 0 {
				character.Processor(l, db).MoveCharacter(event.CharacterId, event.X, event.Y, event.Stance)
			} else if event.Stance != 0 {
				character.Processor(l, db).UpdateStance(event.CharacterId, event.Stance)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [CharacterMovementEvent]")
		}
	}
}
