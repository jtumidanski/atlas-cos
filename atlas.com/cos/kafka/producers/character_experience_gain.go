package producers

import (
	"context"
	"log"
)

type characterExperienceGainEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

var CharacterExperienceGain = func(l *log.Logger, ctx context.Context) *characterExperienceGain {
	return &characterExperienceGain{
		l:   l,
		ctx: ctx,
	}
}

type characterExperienceGain struct {
	l   *log.Logger
	ctx context.Context
}

func (e *characterExperienceGain) Emit(characterId uint32, personalGain uint32, partyGain uint32, show bool, chat bool, white bool) {
	event := &characterExperienceGainEvent{characterId, personalGain, partyGain, show, chat, white}
	produceEvent(e.l, "TOPIC_CHARACTER_EXPERIENCE_EVENT", createKey(int(characterId)), event)
}
