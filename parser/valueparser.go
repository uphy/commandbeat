package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	// ValueType represents the type of field value.
	ValueType string
	// ValueParser parses the values.
	ValueParser interface {
		Parse(value string) (interface{}, error)
	}
	stringParser struct {
	}
	intParser struct {
	}
	floatParser struct {
	}
	boolParser struct {
	}
	timestampParser struct {
		format string
	}
)

var (
	defaultStringParser = &stringParser{}
	defaultIntParser    = &intParser{}
	defaultFloatParser  = &floatParser{}
	defaultBoolParser   = &boolParser{}
)

const (
	// ValueTypeString is the constant value represents 'string' value.
	ValueTypeString ValueType = "string"
	// ValueTypeInt is the constant value represents 'int' value.
	ValueTypeInt = "int"
	// ValueTypeFloat is the constant value represents 'float' value.
	ValueTypeFloat = "float"
	// ValueTypeBool is the constant value represents 'bool' value.
	ValueTypeBool = "bool"
	// ValueTypeTimestamp is the constant value represents 'timestamp' value.
	ValueTypeTimestamp = "timestamp"
	// OptionNameFormat is the option key for creating new ValueParser.
	OptionNameFormat = "format"
)

func newValueParser(valueType ValueType, config map[string]interface{}) (ValueParser, error) {
	switch valueType {
	case "":
		fallthrough
	case ValueTypeString:
		return defaultStringParser, nil
	case ValueTypeInt:
		return defaultIntParser, nil
	case ValueTypeFloat:
		return defaultFloatParser, nil
	case ValueTypeBool:
		return defaultBoolParser, nil
	case ValueTypeTimestamp:
		var timestampFieldConfig TimestampFieldConfig
		if err := convert(config, &timestampFieldConfig); err != nil {
			return nil, err
		}
		return &timestampParser{timestampFieldConfig.Format}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", valueType)
	}
}

func (s *stringParser) Parse(value string) (interface{}, error) {
	return value, nil
}

func (i *intParser) Parse(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (i *floatParser) Parse(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

func (i *boolParser) Parse(value string) (interface{}, error) {
	b := strings.ToLower(value)
	switch b {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "yes":
		return true, nil
	case "no":
		return false, nil
	}
	return nil, fmt.Errorf("unsupported bool value: %s", value)
}

func (i *timestampParser) Parse(value string) (interface{}, error) {
	return time.Parse(i.format, value)
}
