package configuration

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type Registry struct {
	c *Configuration
}

var once sync.Once
var registry *Registry

func Get() *Configuration {
	once.Do(func() {
		c, err := loadConfiguration()
		if err != nil {
			log.Fatalf("Retrieving configuration for service. %s.", err.Error())
		}
		registry = &Registry{
			c: c,
		}
	})
	return registry.c
}

type Configuration struct {
	UseStarting4Ap          bool    `yaml:"useStarting4Ap"`
	UseAutoAssignStartersAp bool    `yaml:"useAutoAssignStartersAp"`
	ExpSplitCommonMod       float32 `yaml:"expSplitCommonMod"`
	ExpSplitMvpMod          float32 `yaml:"expSplitMvpMod"`
	MaxAp                   uint16  `yaml:"maxAp"`
	UseRandomizeHpMpGain    bool    `yaml:"useRandomizeHpMpGain"`
	UseEnforceJobSpRange    bool    `yaml:"useEnforceJobSpRange"`
}
