package query

type StreamEventType string

const (
	StreamEventStarted StreamEventType = "started"
	StreamEventDelta   StreamEventType = "delta"
	StreamEventStep    StreamEventType = "step"
	StreamEventDone    StreamEventType = "done"
	StreamEventError   StreamEventType = "error"
)

type QueryStream <-chan StreamResult

type StreamResult struct {
	Type     StreamEventType
	RunID    string
	Content  string
	Step     *StepResult
	Done     bool
	Err      error
	Metadata map[string]interface{}
}
