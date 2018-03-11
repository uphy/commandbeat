package parser

import (
	"fmt"
	"testing"
)

func newField(name string, valueType ValueType, options map[string]interface{}) (*Field, error) {
	m := map[string]interface{}{}
	for k, v := range options {
		m[k] = v
	}
	m["name"] = name
	m["type"] = valueType
	return newFieldFromConfig(m)
}

func TestCSVParser_Parse(t *testing.T) {
	field1, _ := newField("field1", ValueTypeString, nil)
	field2, _ := newField("field2", ValueTypeInt, nil)
	field3, _ := newField("field3", ValueTypeFloat, nil)
	parser, _ := NewCSVParser(field1, field2, field3)
	m, err := parser.Parse("a,1,0.5")
	fmt.Printf("%#v\n", field2)
	if err != nil {
		t.Errorf("csv cannot be parsed: %s", err)
	}
	if m["field1"] != "a" {
		t.Error("incorrect value: field1")
	}
	if m["field2"] != int64(1) {
		t.Error("incorrect value: field2")
	}
	if m["field3"] != float64(0.5) {
		t.Error("incorrect value: field3")
	}
	if _, err := parser.Parse("a"); err == nil {
		t.Error("inconsistent size not detected.")
	}
	if _, err := parser.Parse("a,a,a"); err == nil {
		t.Error("parse error not thrown.")
	}
	if _, err := NewCSVParser(field1, field1); err == nil {
		t.Error("duplicated field incorrectly allowed.")
	}
	if _, err := newField("a", "unsupportedtype", nil); err == nil {
		t.Error("unsupported type accepted")
	}
}

func TestCSVParserFactory_Create(t *testing.T) {
	factory := &csvFactory{}
	parser, err := factory.Create(map[string]interface{}{
		"fields": []map[string]interface{}{
			{
				"type": "string",
				"name": "field1",
			},
			{
				"type": "int",
				"name": "field2",
			},
		},
	})
	if err != nil {
		t.Errorf("cannot create csv parser factory: %s", err)
	}
	if parser == nil {
		t.Error("parser == nil")
	}
}

func TestCSVParserFactory_CreateInvalidValue(t *testing.T) {
	factory := &csvFactory{}
	_, err := factory.Create(map[string]interface{}{
		"fields": 1,
	})
	if err == nil {
		t.Errorf("should fail")
	}
}

func TestCSVParserFactory_CreateFieldInvalidValue(t *testing.T) {
	factory := &csvFactory{}
	_, err := factory.Create(map[string]interface{}{
		"fields": []map[string]interface{}{
			{
				"type": "unsupportedtype",
				"name": "field1",
			},
		},
	})
	if err == nil {
		t.Errorf("should fail")
	}
}
