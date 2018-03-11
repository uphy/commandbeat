package parser

type (
	// Field represents a document field for Elasticsearch.
	Field struct {
		Name        string
		valueParser ValueParser
	}
	// FieldConfig is the base struct for unmarshalling field config.
	FieldConfig struct {
		Type ValueType `yaml:"type"`
		Name string    `yaml:"name"`
	}
	// TimestampFieldConfig is a struct for unmarhsalling timestamp field config.
	TimestampFieldConfig struct {
		FieldConfig
		Format string `yaml:"format"`
	}
)

func newFieldFromConfig(config map[string]interface{}) (*Field, error) {
	var fieldConfig FieldConfig
	if err := convert(config, &fieldConfig); err != nil {
		return nil, err
	}
	var valueParser ValueParser
	valueParser, err := newValueParser(fieldConfig.Type, config)
	if err != nil {
		return nil, err
	}
	return &Field{fieldConfig.Name, valueParser}, nil
}
