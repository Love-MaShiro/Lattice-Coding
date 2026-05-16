package application

import (
	"context"
	"fmt"
)

type WorkflowState struct {
	RunID      string
	WorkingDir string
	Input      map[string]interface{}
	Values     map[string]interface{}
	Output     string
}

type DAGNode struct {
	Key string
	Run func(ctx context.Context, state *WorkflowState) error
}

type DAGEdge struct {
	From string
	To   string
}

type DAGExecutor struct{}

func NewDAGExecutor() *DAGExecutor {
	return &DAGExecutor{}
}

func (e *DAGExecutor) Execute(ctx context.Context, nodes []DAGNode, edges []DAGEdge, state *WorkflowState) error {
	if state.Values == nil {
		state.Values = map[string]interface{}{}
	}
	ordered, err := topologicalOrder(nodes, edges)
	if err != nil {
		return err
	}
	for _, node := range ordered {
		if node.Run == nil {
			continue
		}
		if err := node.Run(ctx, state); err != nil {
			return fmt.Errorf("%s failed: %w", node.Key, err)
		}
	}
	return nil
}

func topologicalOrder(nodes []DAGNode, edges []DAGEdge) ([]DAGNode, error) {
	nodeByKey := map[string]DAGNode{}
	inDegree := map[string]int{}
	children := map[string][]string{}
	for _, node := range nodes {
		if node.Key == "" {
			return nil, fmt.Errorf("node key is required")
		}
		if _, exists := nodeByKey[node.Key]; exists {
			return nil, fmt.Errorf("duplicated node key: %s", node.Key)
		}
		nodeByKey[node.Key] = node
		inDegree[node.Key] = 0
	}
	for _, edge := range edges {
		if _, ok := nodeByKey[edge.From]; !ok {
			return nil, fmt.Errorf("edge source node not found: %s", edge.From)
		}
		if _, ok := nodeByKey[edge.To]; !ok {
			return nil, fmt.Errorf("edge target node not found: %s", edge.To)
		}
		children[edge.From] = append(children[edge.From], edge.To)
		inDegree[edge.To]++
	}
	queue := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if inDegree[node.Key] == 0 {
			queue = append(queue, node.Key)
		}
	}
	ordered := make([]DAGNode, 0, len(nodes))
	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]
		ordered = append(ordered, nodeByKey[key])
		for _, child := range children[key] {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}
	if len(ordered) != len(nodes) {
		return nil, fmt.Errorf("workflow graph contains cycle")
	}
	return ordered, nil
}
