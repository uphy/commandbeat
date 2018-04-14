package parser

import "testing"

type Foo struct {
	A string
}

func TestConvert(t *testing.T) {
	v := map[string]interface{}{
		"a": "aaaa",
	}
	var foo Foo
	if err := convert(v, &foo); err != nil {
		t.Errorf("convertion failed: %s", err)
	}
	if foo.A != v["a"] {
		t.Errorf("convertion failed")
	}
}

func TestNewParser(t *testing.T) {
	if _, err := NewParser(map[string]interface{}{}); err != nil {
		t.Error(err)
	}
	if _, err := NewParser(map[string]interface{}{
		"type": "csv",
	}); err != nil {
		t.Error(err)
	}
	if _, err := NewParser(map[string]interface{}{
		"type": struct{}{},
	}); err == nil {
		t.Error("should fail")
	}
	if _, err := NewParser(map[string]interface{}{
		"type": "unsupportedtype",
	}); err == nil {
		t.Error("should fail")
	}
}
