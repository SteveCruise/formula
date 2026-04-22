package formula

import (
	"errors"
	"reflect"
	"testing"
)

func TestEngineEval(t *testing.T) {
	engine := NewEngine()
	got, err := engine.Eval("f1 + f2 * 2", map[string]any{
		"f1": 3,
		"f2": 4,
	})
	if err != nil {
		t.Fatalf("Eval() error = %v", err)
	}

	if got != 11 {
		t.Fatalf("Eval() got = %v, want = 11", got)
	}
}

func TestExtractFXVariables(t *testing.T) {
	got, err := ExtractFXVariables("f1 + f20 + f1 + f1f2")
	if err != nil {
		t.Fatalf("ExtractFXVariables() error = %v", err)
	}

	want := []*Join{
		{Id: 1, AssocId: 0},
		{Id: 1, AssocId: 2},
		{Id: 20, AssocId: 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractFXVariables() got = %v, want = %v", got, want)
	}
}

func TestExtractFXVariables_DeduplicateByJoinValue(t *testing.T) {
	got, err := ExtractFXVariables("f1 + f01 + f1f2 + f01f2")
	if err != nil {
		t.Fatalf("ExtractFXVariables() error = %v", err)
	}

	want := []*Join{
		{Id: 1, AssocId: 0},
		{Id: 1, AssocId: 2},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractFXVariables() got = %v, want = %v", got, want)
	}
}

func TestExtractFXVariables_InvalidAny(t *testing.T) {
	_, err := ExtractFXVariables("f1f2f3 + 1")
	if err == nil {
		t.Fatalf("ExtractFXVariables() error = nil, want non-nil")
	}

	var ruleErr *VariableRuleError
	if !errors.As(err, &ruleErr) {
		t.Fatalf("ExtractFXVariables() error = %v, want *VariableRuleError", err)
	}

	if ruleErr.Name != "f1f2f3" && ruleErr.Name != "1" {
		t.Fatalf("VariableRuleError.Name = %q, want %q or %q", ruleErr.Name, "f1f2f3", "1")
	}
}

func TestExtractFXVariables_InvalidNonFx(t *testing.T) {
	_, err := ExtractFXVariables("f1 + abc")
	if err == nil {
		t.Fatalf("ExtractFXVariables() error = nil, want non-nil")
	}

	var ruleErr *VariableRuleError
	if !errors.As(err, &ruleErr) {
		t.Fatalf("ExtractFXVariables() error = %v, want *VariableRuleError", err)
	}
}
