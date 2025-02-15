package rules

import (
	"context"
	"fmt"
)

type (
	Arguments map[string]interface{}

	Condition func(ctx context.Context, args Arguments) bool

	Action func(ctx context.Context, args Arguments) error

	Operator func(a, b bool) bool

	Rule struct {
		Condition Condition
		Action    Action
		Priority  int
	}
)

func Combine(operator Operator, initial bool, conditions ...Condition) Condition {
	return func(ctx context.Context, args Arguments) bool {
		result := initial
		for _, condition := range conditions {
			result = operator(result, condition(ctx, args))
		}
		return result
	}
}

func AllOp(a, b bool) bool {
	return a && b
}

func AnyOp(a, b bool) bool {
	return a || b
}

func NoneOp(a, b bool) bool {
	return a && !b
}

func All(conditions ...Condition) Condition {
	return Combine(AllOp, true, conditions...)
}

func Any(conditions ...Condition) Condition {
	return Combine(AnyOp, false, conditions...)
}

func None(conditions ...Condition) Condition {
	return Combine(NoneOp, true, conditions...)
}

func NewRule(condition Condition, action Action, priority int) *Rule {
	return &Rule{Condition: condition, Action: action, Priority: priority}
}

func (r *Rule) Do(ctx context.Context, args Arguments) error {
	if r.Condition == nil {
		return fmt.Errorf("missing Condition")
	}

	if !r.Condition(ctx, args) {
		return fmt.Errorf("condition was not met")
	}

	if r.Action != nil {
		return r.Action(ctx, args)
	}

	return fmt.Errorf("missing Action")
}
