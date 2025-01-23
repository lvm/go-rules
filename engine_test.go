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
	isEven := func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	}
	isPositive := func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n > 0
	}

	printIsEven := func(c context.Context, args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		fmt.Printf("Success: %d is even!\n", n)
		return nil
	}
	printIsPositive := func(c context.Context, args Arguments) error {
		n, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		fmt.Printf("Success: %d is positive!\n", n)
		return nil
	}

	ctx := context.TODO()
	even := NewRule(isEven, printIsEven, 1)
	positive := NewRule(isPositive, printIsPositive, 2)

	t.Run("AllMatch", func(t *testing.T) {

		ruleEngine := NewRuleEngine(ctx, AllMatch, logger)
		ruleEngine.AddRules(even, positive)

		if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("AnyMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(ctx, AnyMatch, logger)
		ruleEngine.AddRules(even, positive)

		if err := ruleEngine.Execute(Arguments{"number": 3}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("NoneMatch", func(t *testing.T) {
		ruleEngine := NewRuleEngine(ctx, NoneMatch, logger)
		ruleEngine.AddRules(even, positive)

		if err := ruleEngine.Execute(Arguments{"number": -3}); err != nil {
			t.Errorf("Expected nil, but got error (a condition matched)")
		}
	})
}

func TestRuleEngineConditionNotMet(t *testing.T) {
	isEven := func(c context.Context, args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	}
	printErr := func(c context.Context, args Arguments) error {
		_, ok := args["number"].(int)
		if !ok {
			return fmt.Errorf("%+v is not a number", args["number"])
		}
		return fmt.Errorf("this shouldn't run")
	}

	even := NewRule(isEven, printErr, 1)

	ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
	ruleEngine.AddRules(even)

	if err := ruleEngine.Execute(Arguments{"number": 5}); err == nil {
		t.Errorf("Expected an error due to condition not being met, but got nil")
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
		isEven := func(c context.Context, args Arguments) bool {
			if c.Value("ForcePass") != nil {
				return true
			}

			n, _ := args["number"].(int)
			return n%2 == 0
		}

		printIsEven := func(c context.Context, args Arguments) error {
			n, ok := args["number"].(int)
			if !ok {
				return fmt.Errorf("%+v is not a number", args["number"])
			}
			fmt.Printf("Success: %d is even!\n", n)
			return nil
		}

		evenEvenIsNotEven := NewRule(isEven, printIsEven, 1)

		ruleEngine := NewRuleEngine(context.TODO(), AllMatch, logger)
		ruleEngine.AddRules(evenEvenIsNotEven)

		if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		ruleEngine.SetContext("ForcePass", 1)
		if err := ruleEngine.Execute(Arguments{"number": 5}); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

}
