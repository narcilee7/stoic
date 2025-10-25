package planner

import (
	"fmt"
	"sort"
)

// NewRuleSet creates a new RuleSet
func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make([]*Rule, 0),
	}
}

// AddRule adds a new rule to the set
func (rs *RuleSet) AddRule(rule *Rule) {
	rs.rules = append(rs.rules, rule)
}

// RemoveRule removes a rule by name
func (rs *RuleSet) RemoveRule(name string) bool {
	for i, rule := range rs.rules {
		if rule.Name == name {
			rs.rules = append(rs.rules[:i], rs.rules[i+1:]...)
			return true
		}
	}
	return false
}

// GetRule returns a rule by name
func (rs *RuleSet) GetRule(name string) *Rule {
	for _, rule := range rs.rules {
		if rule.Name == name {
			return rule
		}
	}
	return nil
}

// GetAllRules returns a copy of all rules in the set
func (rs *RuleSet) GetAllRules() []*Rule {
	rules := make([]*Rule, len(rs.rules))
	for i, rule := range rs.rules {
		// Create a shallow copy of the rule
		r := *rule
		rules[i] = &r
	}
	return rules
}

// RuleCount returns the number of rules in the set
func (rs *RuleSet) RuleCount() int {
	return len(rs.rules)
}

// Clear removes all rules from the set
func (rs *RuleSet) Clear() {
	rs.rules = make([]*Rule, 0)
}

// Clone creates a deep copy of the RuleSet
func (rs *RuleSet) Clone() *RuleSet {
	newSet := NewRuleSet()
	for _, rule := range rs.rules {
		newRule := *rule
		newSet.AddRule(&newRule)
	}
	return newSet
}

// Match finds all rules that match the given context
func (rs *RuleSet) Match(ctx interface{}) ([]*Rule, error) {
	var matched []*Rule

	for _, rule := range rs.rules {
		matches, err := rule.Condition(ctx)
		if err != nil {
			return nil, fmt.Errorf("error evaluating rule %s: %w", rule.Name, err)
		}
		if matches {
			matched = append(matched, rule)
		}
	}

	// Sort matched rules by priority (highest first)
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Priority > matched[j].Priority
	})

	return matched, nil
}

// MatchFirst returns the first matching rule
func (rs *RuleSet) MatchFirst(ctx interface{}) (*Rule, error) {
	rules, err := rs.Match(ctx)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return nil, nil
	}
	return rules[0], nil
}

// DefaultRules returns a set of default rules
func DefaultRules() *RuleSet {
	rs := NewRuleSet()

	// Example rule: High priority event
	rs.AddRule(&Rule{
		Name:        "high_priority_event",
		Description: "Handle high priority events immediately",
		Priority:    100,
		Condition: func(ctx interface{}) (bool, error) {
			if event, ok := ctx.(map[string]interface{}); ok {
				if priority, ok := event["priority"].(string); ok && priority == "high" {
					return true, nil
				}
			}
			return false, nil
		},
		Action: "handle_immediately",
		Params: map[string]interface{}{
			"timeout": "5s",
		},
	})

	// Example rule: Low priority event
	rs.AddRule(&Rule{
		Name:        "low_priority_event",
		Description: "Handle low priority events in background",
		Priority:    10,
		Condition: func(ctx interface{}) (bool, error) {
			if event, ok := ctx.(map[string]interface{}); ok {
				if priority, ok := event["priority"].(string); ok && priority == "low" {
					return true, nil
				}
			}
			return false, nil
		},
		Action: "queue_for_later",
		Params: map[string]interface{}{
			"queue": "background",
		},
	})

	return rs
}
