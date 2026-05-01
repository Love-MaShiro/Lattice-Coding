import { get, post, put, del } from './request'
import type { Agent, AgentForm } from '@/types/agent'
import type { PageResult, PageQuery } from '@/types/api'

export const agentApi = {
  list: (params?: PageQuery) => get<PageResult<Agent>>('/agents', { params }),

  get: (id: number) => get<Agent>(`/agents/${id}`),

  create: (data: AgentForm) => post<Agent>('/agents', data),

  update: (id: number, data: AgentForm) => put<Agent>(`/agents/${id}`, data),

  delete: (id: number) => del<void>(`/agents/${id}`),

  run: (id: number, input: string) => post<{ runId: string }>(`/agents/${id}/run`, { input })
}
