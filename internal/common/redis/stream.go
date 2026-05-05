package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"lattice-coding/internal/common/errors"
)

type StreamConfig struct {
	StreamName   string
	GroupName    string
	ConsumerName string
	BlockTimeout time.Duration
	Count        int64
	AutoAck      bool
	StartID      string
}

type StreamMessage struct {
	ID      string
	Fields  map[string]interface{}
	Created time.Time
}

func XAdd(ctx context.Context, client *Client, stream string, maxLen int64, values map[string]interface{}) (string, error) {
	valuesCopy := make(map[string]interface{}, len(values))
	for k, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return "", errors.CacheErr("marshal stream value failed: " + err.Error())
		}
		valuesCopy[k] = string(data)
	}

	result, err := client.XAdd(ctx, &goredis.XAddArgs{
		Stream: stream,
		MaxLen: maxLen,
		Approx: false,
		Values: valuesCopy,
	}).Result()
	if err != nil {
		return "", errors.CacheErr("xadd failed: " + err.Error())
	}
	return result, nil
}

func XRead(ctx context.Context, client *Client, cfg StreamConfig) ([]StreamMessage, error) {
	if cfg.BlockTimeout == 0 {
		cfg.BlockTimeout = 5 * time.Second
	}
	if cfg.Count == 0 {
		cfg.Count = 100
	}

	streams := []string{cfg.StreamName}
	if cfg.StartID == "" {
		streams = append(streams, "$")
	} else {
		streams = append(streams, cfg.StartID)
	}

	result, err := client.XRead(ctx, &goredis.XReadArgs{
		Streams: streams,
		Count:   cfg.Count,
		Block:   cfg.BlockTimeout,
	}).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, errors.CacheErr("xread failed: " + err.Error())
	}

	var messages []StreamMessage
	for _, stream := range result {
		for _, msg := range stream.Messages {
			fields := make(map[string]interface{})
			for k, v := range msg.Values {
				if str, ok := v.(string); ok {
					var value interface{}
					if err := json.Unmarshal([]byte(str), &value); err == nil {
						fields[k] = value
					} else {
						fields[k] = v
					}
				} else {
					fields[k] = v
				}
			}
			messages = append(messages, StreamMessage{
				ID:     msg.ID,
				Fields: fields,
			})
		}
	}
	return messages, nil
}

func XReadGroup(ctx context.Context, client *Client, cfg StreamConfig) ([]StreamMessage, error) {
	if cfg.GroupName == "" {
		return nil, errors.CacheErr("group name is required")
	}
	if cfg.ConsumerName == "" {
		return nil, errors.CacheErr("consumer name is required")
	}
	if cfg.BlockTimeout == 0 {
		cfg.BlockTimeout = 5 * time.Second
	}
	if cfg.Count == 0 {
		cfg.Count = 100
	}

	streams := []string{cfg.StreamName, ">"}
	if cfg.StartID != "" {
		streams = []string{cfg.StreamName, cfg.StartID}
	}

	result, err := client.XReadGroup(ctx, &goredis.XReadGroupArgs{
		Group:    cfg.GroupName,
		Consumer: cfg.ConsumerName,
		Streams:  streams,
		Count:    cfg.Count,
		Block:    cfg.BlockTimeout,
	}).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, errors.CacheErr("xreadgroup failed: " + err.Error())
	}

	var messages []StreamMessage
	for _, stream := range result {
		for _, msg := range stream.Messages {
			fields := make(map[string]interface{})
			for k, v := range msg.Values {
				if str, ok := v.(string); ok {
					var value interface{}
					if err := json.Unmarshal([]byte(str), &value); err == nil {
						fields[k] = value
					} else {
						fields[k] = v
					}
				} else {
					fields[k] = v
				}
			}
			messages = append(messages, StreamMessage{
				ID:     msg.ID,
				Fields: fields,
			})
		}
	}
	return messages, nil
}

func XClaim(ctx context.Context, client *Client, stream, group, consumer string, minIdle time.Duration, messageIDs ...string) ([]StreamMessage, error) {
	result, err := client.XClaim(ctx, &goredis.XClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdle,
		Messages: messageIDs,
	}).Result()
	if err != nil {
		return nil, errors.CacheErr("xclaim failed: " + err.Error())
	}

	var messages []StreamMessage
	for _, msg := range result {
		fields := make(map[string]interface{})
		for k, v := range msg.Values {
			fields[k] = v
		}
		messages = append(messages, StreamMessage{
			ID:     msg.ID,
			Fields: fields,
		})
	}
	return messages, nil
}

func XAck(ctx context.Context, client *Client, stream, group string, messageIDs ...string) error {
	if err := client.XAck(ctx, stream, group, messageIDs...).Err(); err != nil {
		return errors.CacheErr("xack failed: " + err.Error())
	}
	return nil
}

func XGroupCreate(ctx context.Context, client *Client, stream, group, startID string) error {
	if err := client.XGroupCreate(ctx, stream, group, startID).Err(); err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return errors.CacheErr("xgroupcreate failed: " + err.Error())
	}
	return nil
}

func XGroupCreateMkStream(ctx context.Context, client *Client, stream, group string) error {
	if err := client.XGroupCreateMkStream(ctx, stream, group, "0").Err(); err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return errors.CacheErr("xgroupcreatemkstream failed: " + err.Error())
	}
	return nil
}

func XPending(ctx context.Context, client *Client, stream, group string) ([]goredis.XPendingExt, error) {
	result, err := client.XPendingExt(ctx, &goredis.XPendingExtArgs{
		Stream: stream,
		Group:  group,
		Start:  "-",
		End:    "+",
		Count:  100,
	}).Result()
	if err != nil {
		return nil, errors.CacheErr("xpending failed: " + err.Error())
	}
	return result, nil
}

func XLen(ctx context.Context, client *Client, stream string) (int64, error) {
	result, err := client.XLen(ctx, stream).Result()
	if err != nil {
		return 0, errors.CacheErr("xlen failed: " + err.Error())
	}
	return result, nil
}

func XDel(ctx context.Context, client *Client, stream string, messageIDs ...string) (int64, error) {
	result, err := client.XDel(ctx, stream, messageIDs...).Result()
	if err != nil {
		return 0, errors.CacheErr("xdel failed: " + err.Error())
	}
	return result, nil
}

func XInfoStream(ctx context.Context, client *Client, stream string) (*goredis.XInfoStream, error) {
	result, err := client.XInfoStream(ctx, stream).Result()
	if err != nil {
		return nil, errors.CacheErr("xinfostream failed: " + err.Error())
	}
	return result, nil
}

func XInfoGroups(ctx context.Context, client *Client, stream string) ([]goredis.XInfoGroup, error) {
	result, err := client.XInfoGroups(ctx, stream).Result()
	if err != nil {
		return nil, errors.CacheErr("xinfogroups failed: " + err.Error())
	}
	return result, nil
}

func PublishEvent(ctx context.Context, client *Client, runID string, eventType string, data interface{}) error {
	streamKey := fmt.Sprintf("stream:run:%s", runID)
	values := map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().UnixMilli(),
		"data":      data,
	}

	_, err := XAdd(ctx, client, streamKey, 1000, values)
	return err
}

func SubscribeRunEvents(ctx context.Context, client *Client, runID, consumerGroup, consumerName string) (<-chan StreamMessage, error) {
	streamKey := fmt.Sprintf("stream:run:%s", runID)

	if err := XGroupCreateMkStream(ctx, client, streamKey, consumerGroup); err != nil {
		return nil, err
	}

	cfg := StreamConfig{
		StreamName:   streamKey,
		GroupName:    consumerGroup,
		ConsumerName: consumerName,
		BlockTimeout: 30 * time.Second,
		Count:        100,
		AutoAck:      true,
	}

	ch := make(chan StreamMessage, 100)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				messages, err := XReadGroup(ctx, client, cfg)
				if err != nil {
					continue
				}

				for _, msg := range messages {
					select {
					case ch <- msg:
						if cfg.AutoAck {
							_ = XAck(ctx, client, streamKey, consumerGroup, msg.ID)
						}
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}
