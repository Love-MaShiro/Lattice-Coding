import { get, post, put, del } from './request'
import type { Agent, AgentDetail, AgentForm } from '@/types/agent'

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

export function buildAgentPayload(form: AgentForm) {
  return {
    name: form.name,
    description: form.description,
    agent_type: 'customer_service',
    model_config_id: form.model_config_id,
    system_prompt: form.system_prompt,
    temperature: form.temperature,
    top_p: 1.0,
    max_tokens: form.max_tokens,
    max_context_turns: form.max_context_turns,
    max_steps: form.max_steps || 20,
    enabled: form.enabled
  }
}

export const agentApi = {
  async list(params?: { page?: number; page_size?: number; keyword?: string }): Promise<TablePage<Agent>> {
    const result = await get<ApiPage<Agent>>('/v1/agents', { params })
    return {
      items: result.data || [],
      total: result.total || 0,
      page: result.page || params?.page || 1,
      size: result.size || params?.page_size || 20
    }
  },

  get: (id: number) => get<{ data: Agent }>(`/v1/agents/${id}`).then((res) => res.data),

  getDetail: (id: number) => get<{ data: AgentDetail }>(`/v1/agents/${id}/detail`).then((res) => res.data),

  create: (data: AgentForm) =>
    post<{ data: Agent }>('/v1/agents', buildAgentPayload(data)).then((res) => res.data),

  update: (id: number, data: AgentForm) =>
    put<{ data: Agent }>(`/v1/agents/${id}`, buildAgentPayload(data)).then((res) => res.data),

  delete: (id: number) => del<void>(`/v1/agents/${id}`),

  enable: (id: number) => post<void>(`/v1/agents/${id}/enable`),

  disable: (id: number) => post<void>(`/v1/agents/${id}/disable`)
}
