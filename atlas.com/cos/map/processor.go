package _map

import "log"

type processor struct {
	l  *log.Logger
}

var Processor = func(l *log.Logger) *processor {
	return &processor{l}
}

func (p processor) FindClosestSpawnPoint(mapId uint32, x int16, y int16) (*Model, error) {

}
