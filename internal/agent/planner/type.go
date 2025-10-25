package planner

import (
	"sync"

	"github.com/stoic/internal/infra"
)

type Rule struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority"` // Higher number means higher priority
	Condition   RuleCondition          `json:"-"`        // Condition to check if the rule applies
	Action      string                 `json:"action"`
	Params      map[string]interface{} `json:"params"`
}

// RuleCondition is a function that evaluates if a rule should be applied
type RuleCondition func(context interface{}) (bool, error)

// RuleSet is a collection of rules with helper methods
type RuleSet struct {
	rules []*Rule
}

// RulesPlanner implements a rule-based planning system
type RulesPlanner struct {
	logger  *infra.Logger
	ruleSet *RuleSet
	mu      sync.RWMutex
}
