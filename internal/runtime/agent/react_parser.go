package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	ReActActionToolCall = "tool_call"
	ReActActionFinal    = "final"
)

type ReActAction struct {
	Type   string                 `json:"type"`
	Reason string                 `json:"reason,omitempty"`
	Tool   string                 `json:"tool,omitempty"`
	Args   map[string]interface{} `json:"args,omitempty"`
	Answer string                 `json:"answer,omitempty"`
}

func ParseReActAction(content string) (*ReActAction, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, errors.New("react action is empty")
	}
	if strings.HasPrefix(trimmed, "```") {
		return nil, errors.New("react action must be raw JSON, not markdown")
	}

	decoder := json.NewDecoder(bytes.NewBufferString(trimmed))
	decoder.DisallowUnknownFields()

	var action ReActAction
	if err := decoder.Decode(&action); err != nil {
		return nil, fmt.Errorf("invalid react JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, errors.New("react action must contain exactly one JSON object")
	}

	switch action.Type {
	case ReActActionToolCall:
		if strings.TrimSpace(action.Tool) == "" {
			return nil, errors.New("tool_call action requires tool")
		}
		if action.Args == nil {
			action.Args = map[string]interface{}{}
		}
	case ReActActionFinal:
		if strings.TrimSpace(action.Answer) == "" {
			return nil, errors.New("final action requires answer")
		}
	default:
		return nil, errors.New("unknown react action type: " + action.Type)
	}

	return &action, nil
}
