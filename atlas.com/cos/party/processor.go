package party

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p *processor) PartyForCharacter(characterId uint32) (*Model, error) {
	return nil, errors.New("party not found")
}