export type SSEEventHandler = (data: any) => void

export interface SSEHandlers {
  onMessageDelta?: SSEEventHandler
  onToolCall?: SSEEventHandler
  onToolResult?: SSEEventHandler
  onRunCompleted?: SSEEventHandler
  onRunFailed?: SSEEventHandler
  onPing?: SSEEventHandler
  onError?: (error: Event) => void
  onClose?: () => void
}

export class RunEventSource {
  private eventSource: EventSource | null = null
  private handlers: SSEHandlers
  private runId: string

  constructor(runId: string, handlers: SSEHandlers) {
    this.runId = runId
    this.handlers = handlers
  }

  connect(): void {
    this.eventSource = new EventSource(`/api/runs/${this.runId}/events`)

    this.eventSource.addEventListener('message.delta', (event: MessageEvent) => {
      this.handlers.onMessageDelta?.(JSON.parse(event.data))
    })

    this.eventSource.addEventListener('tool.call', (event: MessageEvent) => {
      this.handlers.onToolCall?.(JSON.parse(event.data))
    })

    this.eventSource.addEventListener('tool.result', (event: MessageEvent) => {
      this.handlers.onToolResult?.(JSON.parse(event.data))
    })

    this.eventSource.addEventListener('run.completed', (event: MessageEvent) => {
      this.handlers.onRunCompleted?.(JSON.parse(event.data))
    })

    this.eventSource.addEventListener('run.failed', (event: MessageEvent) => {
      this.handlers.onRunFailed?.(JSON.parse(event.data))
    })

    this.eventSource.addEventListener('ping', (event: MessageEvent) => {
      this.handlers.onPing?.(JSON.parse(event.data))
    })

    this.eventSource.onerror = (error: Event) => {
      this.handlers.onError?.(error)
    }
  }

  close(): void {
    if (this.eventSource) {
      this.eventSource.close()
      this.eventSource = null
      this.handlers.onClose?.()
    }
  }
}
