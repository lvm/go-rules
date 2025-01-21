package rules

type Registry struct {
	conditions map[string]Condition
	actions    map[string]Action
}

func NewRegistry() *Registry {
	return &Registry{
		conditions: make(map[string]Condition),
		actions:    make(map[string]Action),
	}
}

func (r *Registry) AddCondition(name string, condition Condition) {
	r.conditions[name] = condition
}

func (r *Registry) GetCondition(name string) (Condition, bool) {
	condition, exists := r.conditions[name]
	return condition, exists
}

func (r *Registry) AddAction(name string, action Action) {
	r.actions[name] = action
}

func (r *Registry) GetAction(name string) (Action, bool) {
	action, exists := r.actions[name]
	return action, exists
}
