import { get, post } from './request'

export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
}

export interface ChatCompletionRequest {
  agentId: number
  messages: ChatMessage[]
}

export interface ChatCompletionResponse {
  content: string
}

export interface ChatStreamResponse {
  delta: string
  done: boolean
}

export const chatApi = {
  completions: (data: ChatCompletionRequest) =>
    post<ChatCompletionResponse>('/chat/completions', data),

  stream: (data: ChatCompletionRequest) =>
    post<ChatStreamResponse>('/chat/stream', data),

  getMessages: (conversationId: string) =>
    get<ChatMessage[]>('/chat/messages', { params: { conversation_id: conversationId } }),

  createMessage: (conversationId: string, message: ChatMessage) =>
    post<void>('/chat/messages', { conversation_id: conversationId, ...message })
}
