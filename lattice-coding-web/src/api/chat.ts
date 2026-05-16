import { del, get, post } from './request'

export type ChatSessionStatus = 'active' | 'archived'
export type ChatMessageRole = 'system' | 'user' | 'assistant' | 'tool'

export interface ChatSession {
  id: number
  title: string
  agent_id: number
  model_config_id: number
  status: ChatSessionStatus | string
  summary: string
  summarized_until_message_id: number
  meta: string
  created_at: string
  updated_at: string
}

export interface ChatMessage {
  id: number
  session_id: number
  role: ChatMessageRole | string
  content: string
  token_count: number
  meta: string
  created_at: string
  updated_at: string
}

export interface CreateChatSessionPayload {
  title?: string
  agent_id: number
}

export interface CreateChatMessagePayload {
  session_id: number
  role?: ChatMessageRole
  content: string
  meta?: string
}

export interface ChatCompletionPayload {
  agent_id?: number
  session_id?: number
  message: string
  mode?: ChatExecutionMode
}

export type ChatExecutionMode = 'direct_chat' | 'fixed_workflow' | 'plan_graph' | 'pure_react'

export interface ChatCompletionResult {
  session_id: number
  message: ChatMessage
  content: string
}

export interface ChatSessionPageQuery {
  page?: number
  page_size?: number
}

export interface ChatMessageListQuery {
  limit?: number
}

export interface ChatStreamHandlers {
  onDelta?: (delta: string) => void
  onEvent?: (event: ChatStreamEvent) => void
  onDone?: (result: ChatCompletionResult) => void
  onError?: (message: string) => void
}

export interface ChatStreamEvent {
  type: string
  run_id?: string
  content?: string
  done?: boolean
  message?: string
  metadata?: Record<string, any>
}

interface ApiPage<T> {
  data: T[]
  total: number
  page: number
  size: number
}

interface TablePage<T> {
  items: T[]
  total: number
  page: number
  size: number
}

function parseSSEBlock(block: string) {
  let event = 'message'
  const data: string[] = []

  block.split('\n').forEach((line) => {
    if (line.startsWith('event:')) {
      event = line.slice(6).trim()
      return
    }
    if (line.startsWith('data:')) {
      data.push(line.slice(5).trimStart())
    }
  })

  return {
    event,
    data: data.join('\n')
  }
}

async function streamComplete(
  data: ChatCompletionPayload,
  handlers: ChatStreamHandlers = {}
): Promise<ChatCompletionResult> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  const token = localStorage.getItem('token')
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch('/api/v1/chat/stream', {
    method: 'POST',
    headers,
    body: JSON.stringify(data)
  })

  if (!response.ok || !response.body) {
    throw new Error(`stream request failed: ${response.status}`)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buffer = ''
  let doneResult: ChatCompletionResult | undefined

  while (true) {
    const { value, done } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const blocks = buffer.split('\n\n')
    buffer = blocks.pop() || ''

    for (const block of blocks) {
      if (!block.trim()) continue
      const parsed = parseSSEBlock(block)
      if (!parsed.data) continue

      const payload = JSON.parse(parsed.data) as ChatStreamEvent
      handlers.onEvent?.(payload)
      if (parsed.event === 'llm.delta') {
        handlers.onDelta?.(payload.content || '')
      } else if (parsed.event === 'run.finished') {
        doneResult = {
          session_id: Number(payload.metadata?.session_id || data.session_id || 0),
          message: {
            id: Number(payload.metadata?.message_id || 0),
            session_id: Number(payload.metadata?.session_id || data.session_id || 0),
            role: 'assistant',
            content: payload.content || '',
            token_count: 0,
            meta: '{}',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          },
          content: payload.content || ''
        }
        handlers.onDone?.(doneResult)
      } else if (parsed.event === 'run.error' || parsed.event === 'error') {
        const message = payload.message || 'stream error'
        handlers.onError?.(message)
        throw new Error(message)
      }
    }
  }

  if (!doneResult) {
    throw new Error('stream closed before done event')
  }
  return doneResult
}

export const chatApi = {
  async listSessions(params?: ChatSessionPageQuery): Promise<TablePage<ChatSession>> {
    const result = await get<ApiPage<ChatSession>>('/v1/chat/sessions', { params })
    return {
      items: result.data || [],
      total: result.total || 0,
      page: result.page || params?.page || 1,
      size: result.size || params?.page_size || 20
    }
  },

  getSession: (id: number) =>
    get<{ data: ChatSession }>(`/v1/chat/sessions/${id}`).then((res) => res.data),

  createSession: (data: CreateChatSessionPayload) =>
    post<{ data: ChatSession }>('/v1/chat/sessions', data).then((res) => res.data),

  deleteSession: (id: number) => del<void>(`/v1/chat/sessions/${id}`),

  compactSession: (id: number) =>
    post<{ data: ChatSession }>(`/v1/chat/sessions/${id}/compact`).then((res) => res.data),

  listMessages: (sessionId: number, params?: ChatMessageListQuery) =>
    get<{ data: ChatMessage[] }>(`/v1/chat/sessions/${sessionId}/messages`, { params }).then(
      (res) => res.data || []
    ),

  createMessage: (data: CreateChatMessagePayload) =>
    post<{ data: ChatMessage }>('/v1/chat/messages', data).then((res) => res.data),

  complete: (data: ChatCompletionPayload) =>
    post<{ data: ChatCompletionResult }>('/v1/chat/completions', data).then((res) => res.data),

  stream: streamComplete
}
