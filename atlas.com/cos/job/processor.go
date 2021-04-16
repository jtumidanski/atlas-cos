package job

import (
	"errors"
	"log"
)

type processor struct {
	l *log.Logger
}

var Processor = func(l *log.Logger) *processor {
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
