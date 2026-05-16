import { get, post, del } from './request'
import type { Run, RunForm, ToolInvocation } from '@/types/run'
import type { PageResult, PageQuery } from '@/types/api'

export const runApi = {
  list: (params?: PageQuery) => get<PageResult<Run>>('/v1/runs', { params }),

  get: (id: string) => get<{ data: Run }>(`/v1/runs/${id}`).then((res) => res.data),

  listToolInvocations: (id: string) =>
    get<{ data: ToolInvocation[] }>(`/v1/runs/${id}/tool-invocations`).then((res) => res.data || []),

  create: (data: RunForm) => post<{ runId: string }>('/v1/runs', data),

  cancel: (id: string) => post<void>(`/v1/runs/${id}/cancel`),

  delete: (id: string) => del<void>(`/v1/runs/${id}`)
}
