package tool

import (
	"context"
	"fmt"
	"unicode/utf8"
)

type ResultProcessor interface {
	Process(ctx context.Context, req ToolRequest, output ToolOutput) (ToolOutput, error)
}

type ResultProcessorFunc func(ctx context.Context, req ToolRequest, output ToolOutput) (ToolOutput, error)

func (f ResultProcessorFunc) Process(ctx context.Context, req ToolRequest, output ToolOutput) (ToolOutput, error) {
	return f(ctx, req, output)
}

type NoopResultProcessor struct{}

func (NoopResultProcessor) Process(_ context.Context, _ ToolRequest, output ToolOutput) (ToolOutput, error) {
	return output, nil
}

type TruncatingResultProcessor struct {
	MaxContentBytes int
}

func NewTruncatingResultProcessor(maxContentBytes int) TruncatingResultProcessor {
	return TruncatingResultProcessor{MaxContentBytes: maxContentBytes}
}

func (p TruncatingResultProcessor) Process(_ context.Context, _ ToolRequest, output ToolOutput) (ToolOutput, error) {
	if p.MaxContentBytes <= 0 || len(output.Content) <= p.MaxContentBytes {
		return output, nil
	}

	originalLen := len(output.Content)
	content := output.Content[:p.MaxContentBytes]
	for !utf8.ValidString(content) && len(content) > 0 {
		content = content[:len(content)-1]
	}

	output.Content = content + "\n[tool output truncated]"
	output.Truncated = true
	if output.Metadata == nil {
		output.Metadata = map[string]interface{}{}
	}
	output.Metadata["truncated"] = true
	output.Metadata["original_content_bytes"] = originalLen
	output.Metadata["max_content_bytes"] = p.MaxContentBytes
	return output, nil
}

type ChainResultProcessor struct {
	processors []ResultProcessor
}

func NewChainResultProcessor(processors ...ResultProcessor) ChainResultProcessor {
	return ChainResultProcessor{processors: processors}
}

func (p ChainResultProcessor) Process(ctx context.Context, req ToolRequest, output ToolOutput) (ToolOutput, error) {
	current := output
	for _, processor := range p.processors {
		if processor == nil {
			continue
		}
		next, err := processor.Process(ctx, req, current)
		if err != nil {
			return ToolOutput{}, err
		}
		current = next
	}
	return current, nil
}

func ErrorContent(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprint(err)
}
