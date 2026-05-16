package prompt

import (
	"context"
	"strings"
)

func (b *PromptBuilder) BuildKnowledgeAnswerPrompt(ctx context.Context, req Request) (*Prompt, error) {
	systemPrompt, err := b.BuildSystemPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	var knowledge strings.Builder
	knowledge.WriteString(systemPrompt)
	knowledge.WriteString("\n\n# Knowledge Answering Mode\n")
	knowledge.WriteString("Answer from the provided knowledge context. If it is insufficient, say what is missing. Do not fabricate citations or facts.\n")
	if req.Knowledge != "" {
		knowledge.WriteString("\n# Knowledge Context\n")
		knowledge.WriteString(req.Knowledge)
	}
	return &Prompt{
		System:   knowledge.String(),
		Messages: []Message{{Role: "user", Content: req.UserInput}},
		Metadata: map[string]interface{}{"kind": "knowledge_answer"},
	}, nil
}
