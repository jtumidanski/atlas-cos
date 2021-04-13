package portal

import (
	"log"
)

type processor struct {
	l  *log.Logger
}

var Processor = func(l *log.Logger) *processor {
	return &processor{l}
}

func (p *processor) GetMapPortalById(mapId uint32, portalId uint32) (*model, error) {

}