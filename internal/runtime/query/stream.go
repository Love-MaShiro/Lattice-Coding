package query

type StreamEventType string

const (
	StreamEventRunStarted  StreamEventType = "run.started"
	StreamEventLLMDelta    StreamEventType = "llm.delta"
	StreamEventLLMDone     StreamEventType = "llm.done"
	StreamEventRunFinished StreamEventType = "run.finished"
	StreamEventRunError    StreamEventType = "run.error"

	StreamEventStarted                 = StreamEventRunStarted
	StreamEventDelta                   = StreamEventLLMDelta
	StreamEventDone                    = StreamEventLLMDone
	StreamEventError                   = StreamEventRunError
	StreamEventStep    StreamEventType = "step"
)

type QueryStream <-chan StreamEvent

type StreamEvent struct {
	Type     StreamEventType
	RunID    string
	Content  string
	Step     *StepResult
	Done     bool
	Err      error
	Metadata map[string]interface{}
}

type StreamResult = StreamEvent
