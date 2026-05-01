import { get, post, del } from './request'
import type { Run, RunForm } from '@/types/run'
import type { PageResult, PageQuery } from '@/types/api'

export const runApi = {
  list: (params?: PageQuery) => get<PageResult<Run>>('/runs', { params }),

  get: (id: number) => get<Run>(`/runs/${id}`),

  create: (data: RunForm) => post<{ runId: string }>('/runs', data),

  cancel: (id: number) => post<void>(`/runs/${id}/cancel`),

  delete: (id: number) => del<void>(`/runs/${id}`)
}
