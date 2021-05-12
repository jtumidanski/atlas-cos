package _map

import (
	"atlas-cos/portal"
	"github.com/sirupsen/logrus"
	"math"
)

func FindClosestSpawnPoint(l logrus.FieldLogger) func(mapId uint32, x int16, y int16) (*portal.Model, error) {
	return func(mapId uint32, x int16, y int16) (*portal.Model, error) {
		portals, err := portal.GetMapPortals(l)(mapId)
		if err != nil {
			return nil, err
		}
		return findClosestSpawnPoint(portals, x, y)
	}
}

func findClosestSpawnPoint(portals []*portal.Model, x int16, y int16) (*portal.Model, error) {
	var closest *portal.Model = nil
	var closestDistance = int16(math.MaxInt16)

	for _, port := range portals {
		if port.IsSpawnPoint() {
			distance := compareDistanceFrom(x, y, port.X(), port.Y())
			if distance < closestDistance {
				closestDistance = distance
				closest = port
			}
		}
	}
	return closest, nil
}

func compareDistanceFrom(referenceX int16, referenceY int16, candidateX int16, candidateY int16) int16 {
	px := referenceX - candidateX
	py := referenceY - candidateY
	return px*px + py*py
}
