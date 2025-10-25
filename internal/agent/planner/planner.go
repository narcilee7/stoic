package planner

type Planner interface {
	Plan(ctx interface{}) (*Plan, error)
}

type Plan struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

func NewPlan(action string, params map[string]interface{}) *Plan {
	if params == nil {
		params = make(map[string]interface{})
	}
	return &Plan{Action: action, Params: params}
}
