package rules

import (
	"context"
	"fmt"
	"testing"
)

func logger(msg string) {
	fmt.Println(msg)
}

func TestRuleEngineMatches(t *testing.T) {

	condIsEven := Condition{Name: "isEven", Fn: func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	}}
	condIsPositive := Condition{Name: "isPositive", Fn: func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n > 0
	}}

	aktPrintIsEven := Action{Name: "printIsEven", Fn: func(c context.Context, args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		fmt.Printf("Success: %d is even!\n", n)
		return nil
	}}
	aktPrintIsPositive := Action{Name: "printIsPositive", Fn: func(c context.Context, args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		fmt.Printf("Success: %d is positive!\n", n)
		return nil
	}}

	ctx := context.TODO()

	t.Run("AllMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(ctx, AllMatch, logger)

		rule1 := NewRule(condIsEven, aktPrintIsEven, 1)
		rule2 := NewRule(condIsPositive, aktPrintIsPositive, 2)

		ruleEngine.AddRules(rule1, rule2)

		if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("AnyMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(ctx, AnyMatch, logger)

		rule1 := NewRule(condIsEven, aktPrintIsEven, 1)
		rule2 := NewRule(condIsPositive, aktPrintIsPositive, 2)

		ruleEngine.AddRules(rule1, rule2)

		if err := ruleEngine.Execute(Arguments{"number": 3}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("NoneMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(ctx, NoneMatch, logger)

		rule1 := NewRule(condIsEven, aktPrintIsEven, 1)
		rule2 := NewRule(condIsPositive, aktPrintIsPositive, 2)

		ruleEngine.AddRules(rule1, rule2)

		if err := ruleEngine.Execute(Arguments{"number": -3}); err != nil {
			t.Errorf("Expected nil, but got error (a condition matched)")
		}
	})
}

func TestRuleEngineConditionNotMet(t *testing.T) {
	condIsEven := Condition{Name: "isEven", Fn: func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	}}

	aktPrintErr := Action{Name: "printErr", Fn: func(c context.Context, args Arguments) error {
		_, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		return fmt.Errorf("this shouldn't run")
	}}

	rule1 := NewRule(condIsEven, aktPrintErr, 1)

	ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
	ruleEngine.AddRules(rule1)

	if err := ruleEngine.Execute(Arguments{"number": 5}); err == nil {
		t.Errorf("Expected an error due to condition not being met, but got nil")
	}
}

func TestRuleEngineInvalidAction(t *testing.T) {
	condIsEven := Condition{Name: "isEven", Fn: func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	}}

	rule := NewRule(condIsEven, Action{Name: "missingInAction"}, 1)

	ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
	ruleEngine.AddRules(rule)

	if err := ruleEngine.Execute(Arguments{"number": 4}); err == nil {
		t.Error("Expected an error for invalid action, but got nil")
	}
}

func TestContext(t *testing.T) {
	ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
	t.Run("Simple Set Get Context", func(t *testing.T) {
		ruleEngine.SetContext("key", "value")
		val := ruleEngine.GetContext("key")

		if val != "value" {
			t.Errorf("Expected context value 'value', got %v", val)
		}
	})

	t.Run("Rules With Context", func(t *testing.T) {
		condIsEven := Condition{Name: "isEven", Fn: func(c context.Context, args Arguments) bool {
			if c.Value("ForcePass") != nil {
				return true
			}

			n, _ := args["number"].(int)
			return n%2 == 0
		}}

		aktPrintIsEven := Action{Name: "printIsEven", Fn: func(c context.Context, args Arguments) error {
			n, ok := args["number"].(int)
			if !ok {
				return fmt.Errorf("%+v is not a number", args["number"])
			}
			fmt.Printf("Success: %d is even!\n", n)
			return nil
		}}

		rule1 := NewRule(condIsEven, aktPrintIsEven, 1)

		ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
		ruleEngine.AddRules(rule1)

		if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		ruleEngine.SetContext("ForcePass", 1)
		if err := ruleEngine.Execute(Arguments{"number": 5}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

}
