package rules

import (
	"fmt"
)

type (
	Arguments map[string]interface{}

	Condition func(args Arguments) bool

	Action func(args Arguments) error

	Operator func(a, b bool) bool

	Rule struct {
		Priority  int       `json:"priority"`
		Condition string    `json:"condition"`
		Args      Arguments `json:"args"`
		Action    string    `json:"action"`
	}
)

func Combine(operator Operator, initial bool, conditions ...Condition) Condition {
	return func(args Arguments) bool {
		result := initial
		for _, condition := range conditions {
			result = operator(result, condition(args))
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

func When(condition string, args Arguments, priority int) *Rule {
	return &Rule{Condition: condition, Args: args, Priority: priority}
}

func (r *Rule) Do(registry *Registry, args Arguments) error {
	conditionFn, ok := registry.GetCondition(r.Condition)
	if !ok {
		return fmt.Errorf("condition not found: %s", r.Condition)
	}

	if !conditionFn(args) {
		return fmt.Errorf("condition '%s' not met", r.Condition)
	}

	actionFn, ok := registry.GetAction(r.Action)
	if !ok {
		return fmt.Errorf("action not found: %s", r.Action)
	}

	return actionFn(args)
}
