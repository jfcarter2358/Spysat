package operation

import (
	"io/ioutil"
	"os"
	"spysat/analyst"
	"spysat/config"
	"spysat/observer"
	"spysat/probe"
	"spysat/window"

	"gopkg.in/yaml.v2"
)

type Operation struct {
	Probes    map[string]probe.Probe       `yaml:"probes"`
	Analysts  map[string]analyst.Analyst   `yaml:"analysts"`
	Window    window.Window                `yaml:"window"`
	Observers map[string]observer.Observer `yaml:"observers"`
}

var Operations Operation

func LoadOperation() error {
	yamlFile, err := ioutil.ReadFile(config.Config.ObserverPath)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &Operations)
	if err != nil {
		return err
	}

	return nil
}

func SaveOperation() error {
	data, err := yaml.Marshal(Operations)
	if err != nil {
		return err
	}

	err = os.WriteFile(config.Config.ObserverPath, data, 0777)
	if err != nil {
		return err
	}

	return nil
}
