package job

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

type processor struct {
	l log.FieldLogger
}

var Processor = func(l log.FieldLogger) *processor {
	return &processor{l}
}

func (p processor) GetByCreateIndex(index uint32) (*Model, error) {
	if index == 0 {
		return &Model{id: Noblesse}, nil
	} else if index == 1 {
		return &Model{id: Beginner}, nil
	} else if index == 2 {
		return &Model{id: Legend}, nil
	}
	return nil, errors.New("job not found for index")
}
