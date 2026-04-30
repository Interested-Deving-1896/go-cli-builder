package resolver

import (
	"os"
	"testing"
)

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