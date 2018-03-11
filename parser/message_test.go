package parser

import "testing"

func TestMessageParser_Parse(t *testing.T) {
	parser := NewMessageParser()
	v, err := parser.Parse("aiueo")
	if err != nil {
		t.Error(err)
	}
	if v["message"] != "aiueo" {
		t.Errorf("expected 'aiueo' but %s", v["message"])
	}
}

func TestMessageParserFactory_Create(t *testing.T) {
	factory := messageFactory{}
	if _, err := factory.Create(nil); err != nil {
		t.Errorf("cannot create json parser factory: %s", err)
	}
}
