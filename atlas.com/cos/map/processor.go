package _map

import (
	"atlas-cos/model"
	"atlas-cos/portal"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math"
)

func FindClosestPortal(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, x int16, y int16) (portal.Model, error) {
	return func(mapId uint32, x int16, y int16) (portal.Model, error) {
		return model.SliceProviderToProviderAdapter[portal.Model](portal.InMapModelProvider(l, span)(mapId), findClosest(x, y))()
	}
}

func FindClosestSpawnPoint(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, x int16, y int16) (portal.Model, error) {
	return func(mapId uint32, x int16, y int16) (portal.Model, error) {
		return model.SliceProviderToProviderAdapter[portal.Model](portal.InMapModelProvider(l, span)(mapId), findClosestSpawnPoint(x, y))()
	}
}

func findClosest(x int16, y int16) model.PreciselyOneFilter[portal.Model] {
	return func(portals []portal.Model) (portal.Model, error) {
		var closest portal.Model
		var closestDistance = int16(math.MaxInt16)
		for _, port := range portals {
			distance := compareDistanceFrom(x, y, port.X(), port.Y())
			if distance < closestDistance {
				closestDistance = distance
				closest = port
			}
		}
		return closest, nil
	}
}

func findClosestSpawnPoint(x int16, y int16) model.PreciselyOneFilter[portal.Model] {
	return func(portals []portal.Model) (portal.Model, error) {
		var closest portal.Model
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
}

func compareDistanceFrom(referenceX int16, referenceY int16, candidateX int16, candidateY int16) int16 {
	px := referenceX - candidateX
	py := referenceY - candidateY
	return px*px + py*py
}
