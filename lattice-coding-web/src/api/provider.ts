import { get, post, put, del } from './request'
import type { ModelConfig, ModelConfigForm, Provider, ProviderForm } from '@/types/provider'
import type { PageResult, PageQuery } from '@/types/api'

export const providerApi = {
  list: (params?: PageQuery) => get<PageResult<Provider>>('/providers', { params }),

  get: (id: number) => get<Provider>(`/providers/${id}`),

  create: (data: ProviderForm) => post<Provider>('/providers', data),

  update: (id: number, data: ProviderForm) => put<Provider>(`/providers/${id}`, data),

  delete: (id: number) => del<void>(`/providers/${id}`),

  test: (id: number) => post<boolean>(`/providers/${id}/test`),

  listModelConfigs: (params?: PageQuery & { provider_id?: number }) =>
    get<PageResult<ModelConfig>>('/model-configs', { params }),

  createModelConfig: (data: ModelConfigForm) => post<ModelConfig>('/model-configs', data),

  updateModelConfig: (id: number, data: Partial<ModelConfigForm>) => put<ModelConfig>(`/model-configs/${id}`, data)
}
