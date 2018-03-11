package parser

import (
	"github.com/elastic/beats/libbeat/common"
)

type (
	// Parser parses stdout as common.MapStr.
	Parser interface {
		Parse(line string) (common.MapStr, error)
	}
)
