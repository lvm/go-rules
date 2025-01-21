package rules

import "testing"

func TestRegistry_ConditionAndActionRegistration(t *testing.T) {
	registry := NewRegistry()
	registry.AddCondition("testCondition", func(args Arguments) bool { return true })
	registry.AddAction("testAction", func(args Arguments) error { return nil })

	if _, exists := registry.conditions["testCondition"]; !exists {
		t.Error("Condition 'testCondition' was not registered")
	}

	if _, exists := registry.actions["testAction"]; !exists {
		t.Error("Action 'testAction' was not registered")
	}
}
