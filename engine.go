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
	Rules         []*Rule         `json:"rules"`
	Context       context.Context `json:"context"`
	ExecutionMode ExecutionMode   `json:"execution_mode"`
	logger        func(string)
	registry      *Registry
}

func NewRuleEngine(mode ExecutionMode, logger func(string), registry *Registry, ctx context.Context) *RuleEngine {
	return &RuleEngine{
		Rules:         []*Rule{},
		Context:       ctx,
		ExecutionMode: mode,
		logger:        logger,
		registry:      registry,
	}
}

func (re *RuleEngine) SetContext(key string, value interface{}) {
	re.Context = context.WithValue(re.Context, key, value)
}

func (re *RuleEngine) GetContext(key string) interface{} {
	return re.Context.Value(key)
}

func (re *RuleEngine) AddRule(rule *Rule) {
	re.Rules = append(re.Rules, rule)
}

func (re *RuleEngine) Execute(args Arguments) error {
	sort.Slice(re.Rules, func(i, j int) bool {
		return re.Rules[i].Priority < re.Rules[j].Priority
	})

	switch re.ExecutionMode {
	case AllMatch:
		for _, rule := range re.Rules {
			if rule.Action != "" {
				if err := rule.Do(re.registry, args); err != nil {
					re.logger(fmt.Sprintf("Rule failed: %v", err))
					return err
				}
			}
		}

	case AnyMatch:
		for _, rule := range re.Rules {
			if rule.Action != "" {
				if err := rule.Do(re.registry, args); err != nil {
					re.logger(fmt.Sprintf("Rule failed: %v", err))
				}
			}
		}

	case NoneMatch:
		failed := true
		for _, rule := range re.Rules {
			if err := rule.Do(re.registry, args); err == nil {
				failed = false
				re.logger(fmt.Sprintf("Rule matched unexpectedly: %v", rule))
			}
		}

		if failed {
			for _, rule := range re.Rules {
				if rule.Action != "" {
					if err := rule.Do(re.registry, args); err != nil {
						re.logger(fmt.Sprintf("Error executing action: %v", err))
					}
				}
			}
			return nil
		}

		return fmt.Errorf("unexpected success, some rules matched")
	}

	return nil
}
