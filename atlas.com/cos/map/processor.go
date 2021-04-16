package _map

import (
	"atlas-cos/portal"
	"log"
)

type processor struct {
	l *log.Logger
}

var Processor = func(l *log.Logger) *processor {
	return &processor{l}
}

func (p processor) FindClosestSpawnPoint(mapId uint32, x int16, y int16) (*portal.Model, error) {
	portals, err := portal.Processor(p.l).GetMapPortals(mapId)
	if err != nil {
		return nil, err
	}
	return p.findClosestSpawnPoint(portals, x, y)
}

func (p processor) findClosestSpawnPoint(portals []*portal.Model, x int16, y int16) (*portal.Model, error) {
	var closest *portal.Model = nil
	var closestDistance int16

	for _, port := range portals {
		if port.IsSpawnPoint() {
			distance := p.compareDistanceFrom(x, y, port.X(), port.Y())
			if distance < closestDistance {
				closestDistance = distance
				closest = port
			}
		}
	}
	return closest, nil
}

func (p processor) compareDistanceFrom(referenceX int16, referenceY int16, candidateX int16, candidateY int16) int16 {
	px := referenceX - candidateX
	py := referenceY - candidateY
	return px*px + py*py
}
