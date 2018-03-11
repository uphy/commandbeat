package parser

import "github.com/elastic/beats/libbeat/common"

type (
	messageParser struct {
	}
	messageFactory struct {
	}
)

func (m *messageFactory) Create(config map[string]interface{}) (Parser, error) {
	return &messageParser{}, nil
}

// NewMessageParser creates message parser.  Message parser use stdout as document field 'message'.
func NewMessageParser() Parser {
	return &messageParser{}
}

func (c *messageParser) Parse(line string) (common.MapStr, error) {
	return common.MapStr{
		"message": line,
	}, nil
}
