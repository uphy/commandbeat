package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/beats/libbeat/common"
)

type (
	csvParser struct {
		fields []*Field
	}
	// CSVConfig is used for unmarshal csv config.
	CSVConfig struct {
		Config
		Fields []map[string]interface{} `yaml:"fields"`
	}
	csvFactory struct {
	}
)

func (c *csvFactory) Create(config map[string]interface{}) (Parser, error) {
	var csvConfig CSVConfig
	if err := convert(config, &csvConfig); err != nil {
		return nil, err
	}
	fields := []*Field{}
	for _, fieldConfig := range csvConfig.Fields {
		field, err := newFieldFromConfig(fieldConfig)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return NewCSVParser(fields...)
}

// NewCSVParser create new CSV Parser.
func NewCSVParser(fields ...*Field) (Parser, error) {
	set := map[string]struct{}{}
	for _, field := range fields {
		if _, exist := set[field.Name]; exist {
			return nil, fmt.Errorf("duplicated field: %s", field.Name)
		}
		set[field.Name] = struct{}{}
	}
	return &csvParser{fields}, nil
}

func (c *csvParser) Parse(line string) (common.MapStr, error) {
	csvReader := csv.NewReader(strings.NewReader(line))
	records, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	if len(records) != len(c.fields) {
		return nil, errors.New("CSV record size differs from the field in config")
	}
	doc := common.MapStr{}
	for i, field := range c.fields {
		parsed, err := field.valueParser.Parse(records[i])
		if err != nil {
			return nil, err
		}
		doc[field.Name] = parsed
	}
	return doc, nil
}
