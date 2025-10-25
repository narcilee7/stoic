package condition

import (
	"context"
	"time"
)

type Condition interface {
	Eval(ctx context.Context, data interface{}) (bool, error)
}

type BasicCondition struct {
	Field    string
	Operator string
	Value    interface{}
}

type AndCondition struct {
	Conditions []Condition
}

type OrCondition struct {
	Conditions []Condition
}

type NotCondition struct {
	Condition Condition
}

type TimeRangeCondition struct {
	Start time.Time
	End   time.Time
}
