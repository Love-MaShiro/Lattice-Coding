package llm

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type EinoModelClient struct {
	factory       *LLMFactory
	modelConfigID uint64
}

func NewEinoModelClient(factory *LLMFactory, modelConfigID uint64) *EinoModelClient {
	return &EinoModelClient{
		factory:       factory,
		modelConfigID: modelConfigID,
	}
}

func ModelConfigClientName(modelConfigID uint64) string {
	return "model_config:" + strconv.FormatUint(modelConfigID, 10)
}

func (c *EinoModelClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	chatModel, err := c.factory.CreateChatModel(ctx, c.modelConfigID)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := chatModel.Generate(ctx, toSchemaMessages(req.Messages), buildOptions(req)...)
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		Content:   strings.TrimSpace(resp.Content),
		LatencyMs: time.Since(start).Milliseconds(),
	}, nil
}

func (c *EinoModelClient) Stream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	chatModel, err := c.factory.CreateChatModel(ctx, c.modelConfigID)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	stream, err := chatModel.Stream(ctx, toSchemaMessages(req.Messages), buildOptions(req)...)
	if err != nil {
		return nil, err
	}

	out := make(chan StreamChunk, 32)
	go func() {
		defer close(out)
		defer stream.Close()

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				out <- StreamChunk{Done: true, LatencyMs: time.Since(start).Milliseconds()}
				return
			}
			if err != nil {
				out <- StreamChunk{Err: err, LatencyMs: time.Since(start).Milliseconds()}
				return
			}
			if chunk == nil || chunk.Content == "" {
				continue
			}
			out <- StreamChunk{
				Content:   chunk.Content,
				LatencyMs: time.Since(start).Milliseconds(),
			}
		}
	}()

	return out, nil
}

func (c *EinoModelClient) Close() error {
	return nil
}

func toSchemaMessages(messages []Message) []*schema.Message {
	result := make([]*schema.Message, 0, len(messages))
	for _, message := range messages {
		result = append(result, &schema.Message{
			Role:    schema.RoleType(message.Role),
			Content: message.Content,
		})
	}
	return result
}

func buildOptions(req ChatRequest) []model.Option {
	opts := make([]model.Option, 0, 4)
	if req.Model != "" {
		opts = append(opts, model.WithModel(req.Model))
	}
	if req.Temperature != nil {
		opts = append(opts, model.WithTemperature(float32(*req.Temperature)))
	}
	if req.TopP != nil {
		opts = append(opts, model.WithTopP(float32(*req.TopP)))
	}
	if req.MaxTokens > 0 {
		opts = append(opts, model.WithMaxTokens(req.MaxTokens))
	}
	return opts
}
