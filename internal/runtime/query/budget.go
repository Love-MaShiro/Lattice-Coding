package query

import (
	"errors"
	"time"
)

type QueryBudget struct {
	MaxSteps     int
	MaxTokens    int
	MaxToolCalls int
	Deadline     time.Time
}

type Budget = QueryBudget

type BudgetTracker struct {
	budget    QueryBudget
	steps     int
	tokens    int
	toolCalls int
}

func NewBudgetTracker(budget QueryBudget) *BudgetTracker {
	return &BudgetTracker{budget: budget}
}

func (t *BudgetTracker) UseStep() error {
	t.steps++
	if t.budget.MaxSteps > 0 && t.steps > t.budget.MaxSteps {
		return ErrBudgetExceeded.WithMessage("query step budget exceeded")
	}
	return t.checkDeadline()
}

func (t *BudgetTracker) UseTokens(tokens int) error {
	if tokens <= 0 {
		return t.checkDeadline()
	}
	t.tokens += tokens
	if t.budget.MaxTokens > 0 && t.tokens > t.budget.MaxTokens {
		return ErrBudgetExceeded.WithMessage("query token budget exceeded")
	}
	return t.checkDeadline()
}

func (t *BudgetTracker) UseToolCall() error {
	t.toolCalls++
	if t.budget.MaxToolCalls > 0 && t.toolCalls > t.budget.MaxToolCalls {
		return ErrBudgetExceeded.WithMessage("query tool budget exceeded")
	}
	return t.checkDeadline()
}

func (t *BudgetTracker) Snapshot() BudgetSnapshot {
	return BudgetSnapshot{
		Steps:     t.steps,
		Tokens:    t.tokens,
		ToolCalls: t.toolCalls,
		Budget:    t.budget,
	}
}

func (t *BudgetTracker) checkDeadline() error {
	if !t.budget.Deadline.IsZero() && time.Now().After(t.budget.Deadline) {
		return errors.Join(ErrBudgetExceeded, ErrQueryTimeout)
	}
	return nil
}

type BudgetSnapshot struct {
	Steps     int
	Tokens    int
	ToolCalls int
	Budget    QueryBudget
}
