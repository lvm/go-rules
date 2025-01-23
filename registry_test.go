package rules

import "testing"

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	t.Run("NewRegistry", func(t *testing.T) {
		if registry == nil {
			t.Error("Registry not created")
		}
	})

	t.Run("AddEngine", func(t *testing.T) {
		registry.AddEngine("test", RuleEngine{})

		if _, exists := registry.engines["test"]; !exists {
			t.Error("RuleEngine 'test' does not exists")
		}
	})

	t.Run("GetEngine", func(t *testing.T) {
		registry.AddEngine("test", RuleEngine{})

		if registry.GetEngine("test") == nil {
			t.Error("RuleEngine 'test' does not exists")
		}
	})

}
