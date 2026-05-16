package tool

import "sync"

type ReadState struct {
	Files map[string]ReadFileState `json:"files,omitempty"`
}

type ReadFileState struct {
	Path       string `json:"path"`
	Version    string `json:"version,omitempty"`
	MTime      int64  `json:"mtime,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
	LastReadAt int64  `json:"last_read_at,omitempty"`
}

type FileReadStateManager interface {
	MarkRead(state ReadFileState)
	Get(path string) (ReadFileState, bool)
	Snapshot() ReadState
}

func NewReadState() ReadState {
	return ReadState{Files: map[string]ReadFileState{}}
}

func (s *ReadState) MarkRead(state ReadFileState) {
	if s.Files == nil {
		s.Files = map[string]ReadFileState{}
	}
	s.Files[state.Path] = state
}

func (s ReadState) Get(path string) (ReadFileState, bool) {
	state, ok := s.Files[path]
	return state, ok
}

type InMemoryFileReadStateManager struct {
	mu    sync.RWMutex
	state ReadState
}

func NewInMemoryFileReadStateManager() *InMemoryFileReadStateManager {
	return &InMemoryFileReadStateManager{state: NewReadState()}
}

func (m *InMemoryFileReadStateManager) MarkRead(state ReadFileState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.MarkRead(state)
}

func (m *InMemoryFileReadStateManager) Get(path string) (ReadFileState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state.Get(path)
}

func (m *InMemoryFileReadStateManager) Snapshot() ReadState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := NewReadState()
	for path, state := range m.state.Files {
		snapshot.Files[path] = state
	}
	return snapshot
}
