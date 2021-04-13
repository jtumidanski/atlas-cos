package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Registry struct {
	c *Configuration
	e error
}

type Configuration struct {
	useStarting4Ap          bool    `yaml:"useStarting4Ap"`
	useAutoAssignStartersAp bool    `yaml:"useAutoAssignStartersAp"`
	expSplitCommonMod       float32 `yaml:"expSplitCommonMod"`
	expSplitMvpMod          float32 `yaml:"expSplitMvpMod"`
	maxAp                   uint16  `yaml:"maxAp"`
	useRandomizeHpMpGain    bool    `yaml:"useRandomizeHpMpGain"`
	useEnforceJobSpRange    bool    `yaml:"useEnforceJobSpRange"`
}

func (c *Configuration) UseRandomizeHPMPGain() bool {
	return c.useRandomizeHpMpGain
}

func (c *Configuration) UseStarting4AP() bool {
	return c.useStarting4Ap
}

func (c *Configuration) MaxAp() uint16 {
	return c.maxAp
}

var configurationRegistryOnce sync.Once
var configurationRegistry *Registry

func Get() (*Configuration, error) {
	configurationRegistryOnce.Do(func() {
		configurationRegistry = &Registry{}
		err := configurationRegistry.loadConfiguration()
		configurationRegistry.e = err
	})
	return configurationRegistry.c, configurationRegistry.e
}

func GetBool(getter func(c *Configuration) bool, def bool) bool {
	c, err := Get()
	if err != nil {
		return def
	}
	return getter(c)
}

func GetUINT16(getter func(c *Configuration) uint16, def uint16) uint16 {
	c, err := Get()
	if err != nil {
		return def
	}
	return getter(c)
}

func (c *Registry) loadConfiguration() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		return err
	}
	c.c = con
	return nil
}
