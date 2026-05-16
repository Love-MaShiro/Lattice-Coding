package query

import "time"

type QueryState struct {
	Request   QueryRequest
	RunID     string
	Mode      ExecutionMode
	StartedAt time.Time
	EndedAt   time.Time
	Budget    *BudgetTracker
	Messages  []Message
	Steps     []StepResult
	Metadata  map[string]interface{}
}

type State = QueryState

func NewState(req QueryRequest) *QueryState {
	runID := req.RunID
	return &QueryState{
		Request:   req,
		RunID:     runID,
		Mode:      req.Mode,
		StartedAt: time.Now(),
		Budget:    NewBudgetTracker(req.Budget),
		Metadata:  map[string]interface{}{},
	}
}

func (s *QueryState) Finish() {
	s.EndedAt = time.Now()
}

type StepResult struct {
	Iteration   int
	Name        string
	Content     string
	IsError     bool
	StartedAt   time.Time
	CompletedAt time.Time
	Metadata    map[string]interface{}
}
