export type SSEEventType =
  | 'message.delta'
  | 'tool.call'
  | 'tool.result'
  | 'run.completed'
  | 'run.failed'
  | 'ping'

export interface SSEMessage {
  event: SSEEventType
  data: any
  runId?: string
}

export interface SSEDelta {
  content: string
}

export interface SSEToolCall {
  tool: string
  input: any
}

export interface SSEToolResult {
  tool: string
  output: any
}

export interface SSERunCompleted {
  runId: string
  output: string
  summary: string
}

export interface SSERunFailed {
  runId: string
  error: string
}
