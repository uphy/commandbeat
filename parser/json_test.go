package parser

import (
	"testing"
)

func TestJSONParser_Parse(t *testing.T) {
	parser := newJSONParser()
	v, err := parser.Parse(`{"a":"aaa","b":1}`)
	if err != nil {
		t.Error(err)
	}
	if v["a"] != "aaa" {
		t.Error("unexpected value: a")
	}
	if v["b"] != float64(1) {
		t.Error("unexpected value: b")
	}
	if _, err := parser.Parse("{aaa"); err == nil {
		t.Error("json error not returned")
	}
}

func TestJSONParserFactory_Create(t *testing.T) {
	factory := jsonFactory{}
	if _, err := factory.Create(nil); err != nil {
		t.Errorf("cannot create json parser factory: %s", err)
	}
}
