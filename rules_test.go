package rules

import (
	"context"
	"fmt"
	"testing"
)

func TestOps(t *testing.T) {
	ctx := context.TODO()
	args := Arguments{}
	truthy := func(c context.Context, args Arguments) bool { return true }
	truthy2 := func(c context.Context, args Arguments) bool { return true }
	falsy := func(c context.Context, args Arguments) bool { return false }

	t.Run("Combine All", func(t *testing.T) {
		allOp := Combine(AllOp, true, truthy, truthy2)
		if !allOp(ctx, args) {
			t.Errorf("Expected combined condition to return true, but got false")
		}

		allOp = Combine(AllOp, true, truthy, falsy)
		if allOp(ctx, args) {
			t.Errorf("Expected combined condition to return false, but got true")
		}
	})

	t.Run("Combine Any", func(t *testing.T) {
		anyOp := Combine(AnyOp, false, truthy, truthy2)
		if !anyOp(ctx, args) {
			t.Errorf("Expected combined condition to return true, but got false")
		}

		anyOp = Combine(AnyOp, false, truthy, falsy)
		if !anyOp(ctx, args) {
			t.Errorf("Expected combined condition to return false, but got true")
		}
	})

	t.Run("Combine None", func(t *testing.T) {
		nonOp := Combine(NoneOp, true, truthy, truthy2)

		if nonOp(ctx, args) {
			t.Errorf("Expected combined condition to return false, but got true")
		}

		truthy3 := func(c context.Context, args Arguments) bool { return true }
		nonOp = Combine(NoneOp, true, truthy, truthy2, truthy3)
		if nonOp(ctx, args) {
			t.Errorf("Expected combined condition to return false, but got true")
		}
	})
}

func TestCombinations(t *testing.T) {
	ctx := context.TODO()
	args := Arguments{}
	truthy := func(c context.Context, args Arguments) bool { return true }
	truthy2 := func(c context.Context, args Arguments) bool { return true }
	falsy := func(c context.Context, args Arguments) bool { return false }

	t.Run("All", func(t *testing.T) {
		allCond := All(truthy, truthy2)
		if !allCond(ctx, args) {
			t.Errorf("Expected All conditions to return true, but got false")
		}

		allCond = All(truthy, falsy)
		if allCond(ctx, args) {
			t.Errorf("Expected All conditions to return false, but got true")
		}
	})

	t.Run("Any", func(t *testing.T) {
		anyCond := Any(truthy, truthy2)
		if !anyCond(ctx, args) {
			t.Errorf("Expected Any conditions to return true, but got false")
		}

		falsy2 := func(c context.Context, args Arguments) bool { return false }
		anyCond = Any(falsy, falsy2)
		if anyCond(ctx, args) {
			t.Errorf("Expected Any conditions to return false, but got true")
		}
	})

	t.Run("None", func(t *testing.T) {
		noneCond := None(truthy, truthy2)
		if noneCond(ctx, args) {
			t.Errorf("Expected None conditions to return true, but got false")
		}

		noneCond = None(truthy, falsy)
		if noneCond(ctx, args) {
			t.Errorf("Expected None conditions to return false, but got true")
		}
	})
}

func TestRuleDo(t *testing.T) {
	isEven := func(c context.Context, args Arguments) bool {
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

	t.Run("Rule.Do()", func(t *testing.T) {
		rule := NewRule(
			isEven,
			printIsEven,
			1,
		)

		err := rule.Do(context.TODO(), Arguments{"number": 4})
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	})

}
