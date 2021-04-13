package party

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p *processor) PartyForCharacter(characterId uint32) (*Model, error) {
	return nil, errors.New("party not found")
}