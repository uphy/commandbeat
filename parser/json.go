package parser

import (
	"encoding/json"

	"github.com/elastic/beats/libbeat/common"
)

type (
	jsonParser struct {
	}
	jsonFactory struct {
	}
)

func (m *jsonFactory) Create(config Config) (Parser, error) {
	return &jsonParser{}, nil
}

// NewJSONParser creates new JSON Parser.
func newJSONParser() Parser {
	return &jsonParser{}
}

func (c *jsonParser) Parse(line string) (common.MapStr, error) {
	var m common.MapStr
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return nil, err
	}
	return m, nil
}
