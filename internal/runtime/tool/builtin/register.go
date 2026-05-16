package builtin

import runtimetool "lattice-coding/internal/runtime/tool"

func RegisterReadOnlyCodingTools(registry *runtimetool.ToolRegistry, stateManager runtimetool.FileReadStateManager) error {
	if registry == nil {
		return nil
	}
	tools := []runtimetool.Tool{
		NewFileReadTool(stateManager),
		NewFileListTool(),
		NewCodeGrepTool(),
		NewGitDiffTool(),
	}
	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

func RegisterCodingTools(registry *runtimetool.ToolRegistry, stateManager runtimetool.FileReadStateManager) error {
	if err := RegisterReadOnlyCodingTools(registry, stateManager); err != nil {
		return err
	}
	if registry == nil {
		return nil
	}
	for _, tool := range []runtimetool.Tool{
		NewFileEditTool(stateManager),
		NewShellRunTool(),
	} {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}
