package rules

import (
	"context"
	"fmt"
	"testing"
)

func TestAllOp(t *testing.T) {
	fn1 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	fn2 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	allOp := Combine(AllOp, true, fn1, fn2)

	if !allOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return true, but got false")
	}

	fn3 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	allOp = Combine(AllOp, true, fn1, fn3)

	if allOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestAnyOp(t *testing.T) {
	fn1 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	fn2 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	anyOp := Combine(AnyOp, false, fn1, fn2)

	if !anyOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return true, but got false")
	}

	fn3 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	anyOp = Combine(AnyOp, false, fn1, fn3)

	if anyOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestNoneOp(t *testing.T) {
	fn1 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	fn2 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	nonOp := Combine(NoneOp, true, fn1, fn2)

	if nonOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}

	fn3 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	nonOp = Combine(NoneOp, true, fn1, fn3)

	if nonOp.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected combined condition to return false, but got true")
	}
}

func TestAll(t *testing.T) {
	fn1 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	fn2 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}

	allCond := All(fn1, fn2)

	if !allCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected All conditions to return true, but got false")
	}

	fn3 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	allCond = All(fn1, fn3)

	if allCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected All conditions to return false, but got true")
	}
}

func TestAny(t *testing.T) {
	fn1 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	fn2 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}

	anyCond := Any(fn1, fn2)

	if !anyCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected Any conditions to return true, but got false")
	}

	fn3 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	anyCond = Any(fn1, fn3)

	if anyCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected Any conditions to return false, but got true")
	}
}

func TestNone(t *testing.T) {
	fn1 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}
	fn2 := Condition{Name: "falseFn", Fn: func(c context.Context, args Arguments) bool { return false }}

	noneCond := None(fn1, fn2)

	if !noneCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected None conditions to return true, but got false")
	}

	fn3 := Condition{Name: "trueFn", Fn: func(c context.Context, args Arguments) bool { return true }}
	noneCond = None(fn1, fn3)

	if noneCond.Fn(context.TODO(), Arguments{}) {
		t.Errorf("Expected None conditions to return false, but got true")
	}
}

func TestRuleDo(t *testing.T) {

	condIsEven := Condition{Name: "isEven", Fn: func(c context.Context, args Arguments) bool {
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

	t.Run("Valid Condition and Action", func(t *testing.T) {
		rule := NewRule(
			condIsEven,
			aktPrintIsEven,
			1,
		)

		err := rule.Do(context.TODO(), Arguments{"number": 4})
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	})

	t.Run("Valid Condition and Missing Action", func(t *testing.T) {
		rule := NewRule(
			condIsEven,
			Action{Name: "fail"},
			1,
		)

		err := rule.Do(context.TODO(), Arguments{"number": 4})
		if err == nil {
			t.Errorf("Expected err, but got: %v", err)
		}
	})
	t.Run("Missing Condition and Valid Action", func(t *testing.T) {
		rule := NewRule(
			Condition{Name: "fail"},
			aktPrintIsEven,
			1,
		)

		err := rule.Do(context.TODO(), Arguments{"number": 4})
		if err == nil {
			t.Errorf("Expected err, but got: %v", err)
		}
	})

}
