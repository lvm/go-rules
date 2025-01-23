package rules

import (
	"context"
	"fmt"
	"sort"
)

type ExecutionMode int

const (
	AllMatch ExecutionMode = iota
	AnyMatch
	NoneMatch
)

type RuleEngine struct {
	Rules         []*Rule
	Context       context.Context
	ExecutionMode ExecutionMode
	logger        func(string)
}

func NewRuleEngine(ctx context.Context, mode ExecutionMode, logger func(string)) *RuleEngine {
	return &RuleEngine{
		Rules:         []*Rule{},
		Context:       ctx,
		ExecutionMode: mode,
		logger:        logger,
	}
}

func (re *RuleEngine) SetContext(key string, value interface{}) {
	re.Context = context.WithValue(re.Context, key, value)
}

func (re *RuleEngine) GetContext(key string) interface{} {
	return re.Context.Value(key)
}

func (re *RuleEngine) AddRules(rules ...*Rule) {
	re.Rules = append(re.Rules, rules...)
}

func (re *RuleEngine) Execute(args Arguments) error {
	sort.Slice(re.Rules, func(i, j int) bool {
		return re.Rules[i].Priority < re.Rules[j].Priority
	})

	switch re.ExecutionMode {
	case AllMatch:
		for _, rule := range re.Rules {
			if err := rule.Do(re.Context, args); err != nil {
				re.logger(fmt.Sprintf("Rule failed: %v", err))
				return err
			}
		}

	case AnyMatch:
		for _, rule := range re.Rules {
			if err := rule.Do(re.Context, args); err == nil {
				re.logger("Rule matched successfully")
				return nil
			}
		}
		re.logger("No rules matched")
		return fmt.Errorf("no rules matched")

	case NoneMatch:
		for _, rule := range re.Rules {
			if err := rule.Do(re.Context, args); err == nil {
				re.logger(fmt.Sprintf("Rule matched unexpectedly: %v", rule))
				return fmt.Errorf("unexpected rule match")
			}
		}

		re.logger("No rules matched, executing all actions")
		for _, rule := range re.Rules {
			if err := rule.Do(re.Context, args); err != nil {
				re.logger(fmt.Sprintf("Error executing action: %v", err))
			}
		}
	}

	return nil
}
