package rules

import (
	"testing"
)

func TestAllOp(t *testing.T) {
	fn1 := func(args Arguments) bool { return true }
	fn2 := func(args Arguments) bool { return true }
	allOp := Combine(AllOp, true, fn1, fn2)

	if !allOp(Arguments{}) {
		t.Errorf("Expected combined condition to return true, but got false")
	}

	fn3 := func(args Arguments) bool { return false }
	allOp = Combine(AllOp, true, fn1, fn3)

	if allOp(Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestAnyOp(t *testing.T) {
	fn1 := func(args Arguments) bool { return false }
	fn2 := func(args Arguments) bool { return true }
	anyOp := Combine(AnyOp, false, fn1, fn2)

	if !anyOp(Arguments{}) {
		t.Errorf("Expected combined condition to return true, but got false")
	}

	fn3 := func(args Arguments) bool { return false }
	anyOp = Combine(AnyOp, false, fn1, fn3)

	if anyOp(Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestNoneOp(t *testing.T) {
	fn1 := func(args Arguments) bool { return true }
	fn2 := func(args Arguments) bool { return true }
	nonOp := Combine(NoneOp, true, fn1, fn2)

	if nonOp(Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}

	fn3 := func(args Arguments) bool { return true }
	nonOp = Combine(NoneOp, true, fn1, fn3)

	if nonOp(Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestAll(t *testing.T) {
	fn1 := func(args Arguments) bool { return true }
	fn2 := func(args Arguments) bool { return true }

	allCond := All(fn1, fn2)

	if !allCond(Arguments{}) {
		t.Errorf("Expected All conditions to return true, but got false")
	}

	fn3 := func(args Arguments) bool { return false }
	allCond = All(fn1, fn3)

	if allCond(Arguments{}) {
		t.Errorf("Expected All conditions to return false, but got true")
	}
}

func TestAny(t *testing.T) {
	fn1 := func(args Arguments) bool { return false }
	fn2 := func(args Arguments) bool { return true }

	anyCond := Any(fn1, fn2)

	if !anyCond(Arguments{}) {
		t.Errorf("Expected Any conditions to return true, but got false")
	}

	fn3 := func(args Arguments) bool { return false }
	anyCond = Any(fn1, fn3)

	if anyCond(Arguments{}) {
		t.Errorf("Expected Any conditions to return false, but got true")
	}
}

func TestNone(t *testing.T) {
	fn1 := func(args Arguments) bool { return false }
	fn2 := func(args Arguments) bool { return false }

	noneCond := None(fn1, fn2)

	if !noneCond(Arguments{}) {
		t.Errorf("Expected None conditions to return true, but got false")
	}

	fn3 := func(args Arguments) bool { return true }
	noneCond = None(fn1, fn3)

	if noneCond(Arguments{}) {
		t.Errorf("Expected None conditions to return false, but got true")
	}
}

func TestRule_Do(t *testing.T) {
	registry := NewRegistry()

	registry.AddCondition("isEven", func(args Arguments) bool {
		n, _ := args["number"].(int)
		return n%2 == 0
	})
	registry.AddAction("increment", func(args Arguments) error {
		return nil
	})

	rule := When("isEven", Arguments{"number": 4}, 1)
	rule.Action = "increment"

	err := rule.Do(registry, Arguments{"number": 4})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	rule.Condition = "nonexistentCondition"
	err = rule.Do(registry, Arguments{"number": 4})
	if err == nil || err.Error() != "condition not found: nonexistentCondition" {
		t.Errorf("Expected 'condition not found' error, but got: %v", err)
	}

	rule.Condition = "isEven"
	rule.Action = "nonexistentAction"
	err = rule.Do(registry, Arguments{"number": 4})
	if err == nil || err.Error() != "action not found: nonexistentAction" {
		t.Errorf("Expected 'action not found' error, but got: %v", err)
	}

	rule.Action = "increment"
	err = rule.Do(registry, Arguments{"number": 5})
	if err == nil || err.Error() != "condition 'isEven' not met" {
		t.Errorf("Expected 'condition not met' error, but got: %v", err)
	}
}
