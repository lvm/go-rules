package rules

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	ruleEngine := NewRuleEngine(AllMatch, func(msg string) {}, NewRegistry(), context.TODO())

	ruleEngine.SetContext("key", "value")
	val := ruleEngine.GetContext("key")

	if val != "value" {
		t.Errorf("Expected context value 'value', got %v", val)
	}
}
func TestRuleEngineMatches(t *testing.T) {
	registry := NewRegistry()

	registry.AddCondition("isEven", func(args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	})

	registry.AddCondition("isPositive", func(args Arguments) bool {
		n, _ := args["number"].(int)
		return n > 0
	})

	registry.AddAction("printIsEven", func(args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}

		fmt.Printf("Success: %d is even!\n", n)
		return nil
	})

	registry.AddAction("printIsPositive", func(args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}

		fmt.Printf("Success: %d is positive!\n", n)
		return nil
	})

	t.Run("AllMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(AllMatch, func(msg string) {}, registry, context.TODO())

		rule1 := When("isEven", Arguments{"number": 4}, 1)
		rule1.Action = "printIsEven"
		ruleEngine.AddRule(rule1)

		rule2 := When("isPositive", Arguments{"number": 4}, 2)
		rule2.Action = "printIsPositive"
		ruleEngine.AddRule(rule2)

		if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("AnyMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(AnyMatch, func(msg string) {}, registry, context.TODO())

		rule1 := When("isEven", Arguments{"number": 3}, 1)
		rule1.Action = "printIsEven"
		ruleEngine.AddRule(rule1)

		rule2 := When("isPositive", Arguments{"number": 3}, 2)
		rule2.Action = "printIsPositive"
		ruleEngine.AddRule(rule2)

		if err := ruleEngine.Execute(Arguments{"number": 3}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("NoneMatch", func(t *testing.T) {

		ruleEngine := NewRuleEngine(NoneMatch, func(msg string) {}, registry, context.TODO())

		rule1 := When("isEven", Arguments{"number": -3}, 1)
		rule1.Action = "printIsEven"
		ruleEngine.AddRule(rule1)

		rule2 := When("isPositive", Arguments{"number": -3}, 2)
		rule2.Action = "printIsPositive"
		ruleEngine.AddRule(rule2)

		if err := ruleEngine.Execute(Arguments{"number": -3}); err != nil {
			t.Errorf("Expected nil, but got error (a condition matched)")
		}
	})
}

func TestRuleEngineConditionNotMet(t *testing.T) {
	registry := NewRegistry()

	registry.AddCondition("isEven", func(args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	})

	registry.AddAction("printErr", func(args Arguments) error {
		_, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}

		return fmt.Errorf("this shouldn't run")
	})

	ruleEngine := NewRuleEngine(AllMatch, func(msg string) {}, registry, context.TODO())

	rule1 := When("isEven", Arguments{"number": 5}, 1)
	rule1.Action = "printErr"
	ruleEngine.AddRule(rule1)

	if err := ruleEngine.Execute(Arguments{"number": 5}); err == nil {
		t.Errorf("Expected an error due to condition not being met, but got err: %v", err)
	}
}

func TestRuleEngineInvalidAction(t *testing.T) {
	registry := NewRegistry()

	registry.AddCondition("isEven", func(args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	})

	ruleEngine := NewRuleEngine(AllMatch, func(msg string) {}, registry, context.TODO())

	rule1 := When("isEven", Arguments{"number": 4}, 1)
	rule1.Action = "missingInAction"
	ruleEngine.AddRule(rule1)

	if err := ruleEngine.Execute(Arguments{"number": 4}); err == nil {
		t.Error("Expected an error for invalid action, but got nil")
	}
}
