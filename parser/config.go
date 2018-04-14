package parser

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type (
	// Factory is the factory interface to create Parser.
	Factory interface {
		Create(config Config) (Parser, error)
	}
	// Config is the parser config.
	Config map[string]interface{}
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
	t, ok := config["type"]
	if !ok {
		return NewDefaultParser(), nil
	}
	tstr, ok := t.(string)
	if !ok {
		return nil, errors.New("parser: invalid type value")
	}
	factory := factories[tstr]
	if factory == nil {
		return nil, fmt.Errorf("unsupported parser: %s", t)
	}
	return factory.Create(Config(config))
}
