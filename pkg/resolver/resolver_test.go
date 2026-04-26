package resolver

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestBindValue_String(t *testing.T) {
	var s string
	val := reflect.ValueOf(&s).Elem()
	if err := BindValue(val, "hello"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if s != "hello" {
		t.Errorf("got %q, want %q", s, "hello")
	}
}

func TestBindValue_Int(t *testing.T) {
	var v int
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "42"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 42 {
		t.Errorf("got %d, want %d", v, 42)
	}
}

func TestBindValue_Int8(t *testing.T) {
	var v int8
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "8"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 8 {
		t.Errorf("got %d, want %d", v, 8)
	}
}

func TestBindValue_Int64(t *testing.T) {
	var v int64
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "9223372036854775807"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 9223372036854775807 {
		t.Errorf("got %d, want %d", v, 9223372036854775807)
	}
}

func TestBindValue_Float32(t *testing.T) {
	var v float32
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "3.14"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 3.14 {
		t.Errorf("got %f, want %f", v, 3.14)
	}
}

func TestBindValue_Float64(t *testing.T) {
	var v float64
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "2.718"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 2.718 {
		t.Errorf("got %f, want %f", v, 2.718)
	}
}

func TestBindValue_BoolTrue(t *testing.T) {
	var v bool
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "true"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if !v {
		t.Error("expected true")
	}
}

func TestBindValue_BoolFalse(t *testing.T) {
	var v bool
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "false"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v {
		t.Error("expected false")
	}
}

func TestBindValue_Duration(t *testing.T) {
	var v time.Duration
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "5s"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if v != 5*time.Second {
		t.Errorf("got %v, want %v", v, 5*time.Second)
	}
}

func TestBindValue_StringSlice(t *testing.T) {
	v := []string{}
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "a"); err != nil {
		t.Fatalf("BindValue failed: %v", err)
	}
	if len(v) != 1 || v[0] != "a" {
		t.Errorf("got %v, want [a]", v)
	}
}

func TestBindValue_IntParseError(t *testing.T) {
	var v int
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "not-a-number"); err == nil {
		t.Error("expected error for invalid integer")
	}
}

func TestBindValue_FloatParseError(t *testing.T) {
	var v float64
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "not-float"); err == nil {
		t.Error("expected error for invalid float")
	}
}

func TestBindValue_BoolParseError(t *testing.T) {
	var v bool
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "maybe"); err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestBindValue_UnsupportedType(t *testing.T) {
	var v complex64
	val := reflect.ValueOf(&v).Elem()
	if err := BindValue(val, "5"); err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestParseBool_TrueVariants(t *testing.T) {
	cases := []string{"1", "t", "T", "true", "TRUE", "True", "yes", "Yes", "YES", "on", "ON"}
	for _, c := range cases {
		got, err := ParseBool(c)
		if err != nil {
			t.Errorf("ParseBool(%q) unexpected error: %v", c, err)
		}
		if !got {
			t.Errorf("ParseBool(%q) = false, want true", c)
		}
	}
}

func TestParseBool_FalseVariants(t *testing.T) {
	cases := []string{"0", "f", "F", "false", "FALSE", "False", "no", "No", "NO", "off", "OFF"}
	for _, c := range cases {
		got, err := ParseBool(c)
		if err != nil {
			t.Errorf("ParseBool(%q) unexpected error: %v", c, err)
		}
		if got {
			t.Errorf("ParseBool(%q) = true, want false", c)
		}
	}
}

func TestParseBool_Invalid(t *testing.T) {
	_, err := ParseBool("invalid")
	if err == nil {
		t.Error("expected error for invalid bool string")
	}
}

func TestGetValue_CLI(t *testing.T) {
	cli := "cli-value"
	val, err := GetValue(&cli, "ENV", "default", false)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if val != "cli-value" {
		t.Errorf("got %q, want %q", val, "cli-value")
	}
}

func TestGetValue_Env(t *testing.T) {
	os.Setenv("TEST_VAR", "env-value")
	defer os.Unsetenv("TEST_VAR")

	val, err := GetValue(nil, "TEST_VAR", "default", false)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if val != "env-value" {
		t.Errorf("got %q, want %q", val, "env-value")
	}
}

func TestGetValue_Default(t *testing.T) {
	val, err := GetValue(nil, "", "default-val", false)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if val != "default-val" {
		t.Errorf("got %q, want %q", val, "default-val")
	}
}

func TestGetValue_Empty(t *testing.T) {
	val, err := GetValue(nil, "", "", false)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if val != "" {
		t.Errorf("got %q, want empty", val)
	}
}

func TestGetValue_Priority(t *testing.T) {
	os.Setenv("PRIORITY_VAR", "env-value")
	defer os.Unsetenv("PRIORITY_VAR")

	cli := "cli-value"
	val, err := GetValue(&cli, "PRIORITY_VAR", "default", false)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if val != "cli-value" {
		t.Errorf("CLI should have priority, got %q", val)
	}
}
