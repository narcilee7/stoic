package planner

import (
	"fmt"

	"github.com/stoic/internal/infra"
)

func NewRulesPlanner(logger *infra.Logger) *RulesPlanner {
	return &RulesPlanner{
		logger:  logger,
		ruleSet: DefaultRules(),
	}
}

func (rp *RulesPlanner) AddRule(rule *Rule) error {
	if rule == nil {
		return fmt.Errorf("rule cannot be nil")
	}
	rp.mu.Lock()
	defer rp.mu.Unlock()

	if existing := rp.ruleSet.GetRule(rule.Name); existing != nil {
		return fmt.Errorf("rule with name '%s' already exists", rule.Name)
	}

	rp.ruleSet.AddRule(rule)
	rp.logger.Debug("Added new rule", "rule", rule.Name)

	return nil
}

func (rp *RulesPlanner) RemoveRule(name string) bool {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	removed := rp.ruleSet.RemoveRule(name)
	if removed {
		rp.logger.Debug("Removed rule", "name", name)
	}
	return removed
}

func (rp *RulesPlanner) GetRule(name string) *Rule {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	return rp.ruleSet.GetRule(name)
}

func (rp *RulesPlanner) RuleCount() int {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	return rp.ruleSet.RuleCount()
}

func (rp *RulesPlanner) Clear() {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.ruleSet.Clear()
	rp.logger.Debug("Cleared all rules")
}

func (rp *RulesPlanner) ReloadRules(rules []*Rule) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	newSet := NewRuleSet()
	for _, rule := range rules {
		newSet.AddRule(rule)
	}

	rp.ruleSet = newSet
	rp.logger.Info("Reloaded rules", "count", newSet.RuleCount())
}

func (rp *RulesPlanner) GetRules() []*Rule {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	return rp.ruleSet.GetAllRules()
}

func (rp *RulesPlanner) Plan(ctx interface{}) (*Plan, error) {
	rule, err := rp.ruleSet.MatchFirst(ctx)
	if err != nil {
		return nil, fmt.Errorf("error matching rules:%w", err)
	}

	if rule == nil {
		rp.logger.Debug("No matching rules found")
		return &Plan{
			Action: "noop",
			Params: map[string]interface{}{"reason": "no_matching_rule"},
		}, nil
	}

	rp.logger.Debug("Matched rule", "rule", rule.Name, "action", rule.Action)
	return &Plan{
		Action: rule.Action,
		Params: rule.Params,
	}, nil
}
