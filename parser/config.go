package parser

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type (
	// Config is base struct for unmarshalling config.
	Config struct {
		Type string `yaml:"type"`
	}
	// Factory is the factory interface to create Parser.
	Factory interface {
		Create(config map[string]interface{}) (Parser, error)
	}
)

var (
	factories = map[string]Factory{}
)

func init() {
	messageFactory := &messageFactory{}
	factories["message"] = messageFactory
	factories[""] = messageFactory
	factories["csv"] = &csvFactory{}
	factories["json"] = &jsonFactory{}
}

func convert(from interface{}, to interface{}) error {
	y, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(y, to)
}

// NewParser creates Parser from config.
func NewParser(config map[string]interface{}) (Parser, error) {
	var c Config
	if err := convert(config, &c); err != nil {
		return nil, err
	}
	factory := factories[c.Type]
	if factory == nil {
		return nil, fmt.Errorf("unsupported parser: %s", c.Type)
	}
	return factory.Create(config)
}
