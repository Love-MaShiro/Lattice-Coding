import { get, post, put, del } from './request'
import type { Provider, ProviderForm } from '@/types/provider'
import type { PageResult, PageQuery } from '@/types/api'

export const providerApi = {
  list: (params?: PageQuery) => get<PageResult<Provider>>('/providers', { params }),

  get: (id: number) => get<Provider>(`/providers/${id}`),

  create: (data: ProviderForm) => post<Provider>('/providers', data),

  update: (id: number, data: ProviderForm) => put<Provider>(`/providers/${id}`, data),

  delete: (id: number) => del<void>(`/providers/${id}`),

  test: (id: number) => post<boolean>(`/providers/${id}/test`)
}
