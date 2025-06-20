package utils

import (
	"testing"

	"github.com/leinonen/go-lisp/pkg/registry"
)

func TestUtilsPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewUtilsPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{
		"frequencies", "group-by", "partition", "interleave", "interpose",
		"flatten", "shuffle", "remove", "keep", "mapcat",
		"take-while", "drop-while", "split-at", "split-with",
		"comp", "partial", "complement", "juxt",
		"union", "intersection", "difference",
	}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestUtilsPlugin_PluginInfo(t *testing.T) {
	plugin := NewUtilsPlugin()

	if plugin.Name() != "utils" {
		t.Errorf("Expected plugin name 'utils', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}
}
