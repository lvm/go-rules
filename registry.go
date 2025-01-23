package rules

type Registry struct {
	engines map[string]RuleEngine
}

func NewRegistry() *Registry {
	return &Registry{engines: make(map[string]RuleEngine)}
}

func (r *Registry) AddEngine(name string, engine RuleEngine) {
	r.engines[name] = engine
}

func (r *Registry) GetEngine(name string) *RuleEngine {
	re, ok := r.engines[name]
	if !ok {
		return nil
	}

	return &re
}
