package parser

import "testing"

func TestNewValueParser(t *testing.T) {
	if _, err := newValueParser(ValueTypeString, nil); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newValueParser("", nil); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newValueParser(ValueTypeInt, nil); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newValueParser(ValueTypeFloat, nil); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newValueParser(ValueTypeBool, nil); err != nil {
		t.Errorf("failed to create field: %s", err)
	}
	if _, err := newValueParser(ValueTypeTimestamp, map[string]interface{}{
		"format": "2006/01/02 15:04:05",
	}); err != nil {
		t.Errorf("failed to create field: %s", err)
	}

	if _, err := newValueParser(ValueTypeTimestamp, map[string]interface{}{
		"format": struct{}{},
	}); err == nil {
		t.Error("expected error")
	}
}

func TestValueParserNew(t *testing.T) {
	for _, typ := range []ValueType{ValueTypeString, ValueTypeInt, ValueTypeFloat} {
		if v, err := newValueParser(typ, nil); v == nil || err != nil {
			t.Errorf("parser == nil or err != nil (type=%v, parser=%v, err=%v)", typ, v, err)
		}
	}
	if v, err := newValueParser("unsupportedtype", nil); v != nil || err == nil {
		t.Errorf("type is not supported type but created new instance (parser=%v, err=%v)", v, err)
	}
	if _, err := newValueParser(ValueTypeTimestamp, map[string]interface{}{
		OptionNameFormat: "2006/01/02 15:04:05",
	}); err != nil {
		t.Errorf("cannot create timestamp ValueParser with format.")
	}
}

func TestValueParser_string(t *testing.T) {
	parser, err := newValueParser(ValueTypeString, nil)
	if err != nil {
		t.Fatal()
	}
	v, err := parser.Parse("aiueo")
	if err != nil {
		t.Errorf("should be parsed. (err=%s)", err)
	}
	if v != "aiueo" {
		t.Errorf("string should be same. (source=%s, parsed=%s)", "aiueo", v)
	}
}

func TestValueParser_int(t *testing.T) {
	parser, err := newValueParser(ValueTypeInt, nil)
	if err != nil {
		t.Fatal()
	}
	v, err := parser.Parse("0")
	if err != nil {
		t.Error(err)
	}
	if v != int64(0) {
		t.Errorf("int parse failed. (source=0, parsed=%d", v)
	}
}

func TestValueParser_float(t *testing.T) {
	parser, err := newValueParser(ValueTypeFloat, nil)
	if err != nil {
		t.Fatal()
	}
	v, err := parser.Parse("0.1")
	if err != nil {
		t.Error(err)
	}
	if v != float64(0.1) {
		t.Errorf("float parse failed. (source=0.1, parsed=%v", v)
	}
}

func TestValueParser_bool(t *testing.T) {
	parser, err := newValueParser(ValueTypeBool, nil)
	if err != nil {
		t.Fatal()
	}
	for _, b := range []string{"true", "yes", "YES", "TRUE"} {
		v, err := parser.Parse(b)
		if err != nil {
			t.Error(err)
		}
		if v != true {
			t.Errorf("bool parse failed. (source=%s, parsed=%v", b, v)
		}
	}
	for _, b := range []string{"false", "no", "NO", "FALSE"} {
		v, err := parser.Parse(b)
		if err != nil {
			t.Error(err)
		}
		if v != false {
			t.Errorf("bool parse failed. (source=%s, parsed=%v", b, v)
		}
	}
	if _, err := parser.Parse("unsupportedbool"); err == nil {
		t.Errorf("shouldn't be parsed.")
	}
}

func TestValueParser_timestamp(t *testing.T) {
	parser, _ := newValueParser(ValueTypeTimestamp, map[string]interface{}{
		OptionNameFormat: "2006/01/02 15:04:05",
	})
	if _, err := parser.Parse("2018/03/11 15:07:53"); err != nil {
		t.Errorf("cannot parse timestamp.")
	}
}
