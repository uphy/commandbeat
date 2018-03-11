package parser

import "testing"

func TestNewFieldFromConfig(t *testing.T) {
	if _, err := newFieldFromConfig(map[string]interface{}{
		"name": "field1",
		"type": "string",
	}); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newFieldFromConfig(map[string]interface{}{
		"name": struct{}{},
	}); err == nil {
		t.Error("expected error")
	}
}
