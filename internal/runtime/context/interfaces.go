package context

import stdcontext "context"

type Builder interface {
	Build(ctx stdcontext.Context, req Request) (*Context, error)
}

type Provider interface {
	Load(ctx stdcontext.Context, req Request) (*Section, error)
}

type WorkspaceGuard interface {
	Resolve(path string) (string, error)
	Contains(path string) bool
}

type Request struct {
	RunID      string
	AgentID    string
	UserID     string
	ProjectID  string
	WorkingDir string
	Metadata   map[string]interface{}
}

type Context struct {
	System      string
	User        string
	Project     string
	Git         GitState
	Workspace   Workspace
	Sections    []Section
	TokenBudget int
}

type Section struct {
	Name     string
	Content  string
	Priority int
	Metadata map[string]interface{}
}

type GitState struct {
	Branch  string
	Commit  string
	IsDirty bool
	Summary string
}

type Workspace struct {
	Root        string
	AllowedDirs []string
}
