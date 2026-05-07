import { get, post, put, del } from './request'
import type {
  ModelConfig,
  ModelConfigForm,
  ModelTestResult,
  Provider,
  ProviderForm,
  ProviderHealthResult,
  ProviderTestResult,
  SyncModelsResult
} from '@/types/provider'
import type { PageResult, PageQuery } from '@/types/api'

export const providerApi = {
  list: (params?: PageQuery & { keyword?: string }) =>
    get<PageResult<Provider>>('/v1/providers', { params }),

  get: (id: number) =>
    get<{ data: Provider }>(`/v1/providers/${id}`).then((res) => res.data),

  create: (data: ProviderForm) =>
    post<{ data: Provider }>('/v1/providers', data).then((res) => res.data),

  update: (id: number, data: Partial<ProviderForm>) =>
    put<{ data: Provider }>(`/v1/providers/${id}`, data).then((res) => res.data),

  delete: (id: number) => del<void>(`/v1/providers/${id}`),

  enable: (id: number) => post<void>(`/v1/providers/${id}/enable`),

  disable: (id: number) => post<void>(`/v1/providers/${id}/disable`),

  test: (id: number) =>
    post<{ data: ProviderTestResult }>(`/v1/providers/${id}/test`).then((res) => res.data),

  healthCheck: (id: number) =>
    post<{ data: ProviderHealthResult }>(`/v1/providers/${id}/health-check`).then((res) => res.data),

  syncModels: (id: number) =>
    post<{ data: SyncModelsResult }>(`/v1/providers/${id}/sync-models`).then((res) => res.data),

  listModelConfigs: (params?: PageQuery & { provider_id?: number }) =>
    get<PageResult<ModelConfig>>('/v1/model-configs', { params }),

  getModelConfig: (id: number) =>
    get<{ data: ModelConfig }>(`/v1/model-configs/${id}`).then((res) => res.data),

  createModelConfig: (data: ModelConfigForm) =>
    post<{ data: ModelConfig }>('/v1/model-configs', data).then((res) => res.data),

  updateModelConfig: (id: number, data: Partial<ModelConfigForm>) =>
    put<{ data: ModelConfig }>(`/v1/model-configs/${id}`, data).then((res) => res.data),

  deleteModelConfig: (id: number) => del<void>(`/v1/model-configs/${id}`),

  enableModelConfig: (id: number) => post<void>(`/v1/model-configs/${id}/enable`),

  disableModelConfig: (id: number) => post<void>(`/v1/model-configs/${id}/disable`),

  setDefaultModelConfig: (id: number) => post<void>(`/v1/model-configs/${id}/default`),

  testModelConfig: (id: number) =>
    post<{ data: ModelTestResult }>(`/v1/model-configs/${id}/test`).then((res) => res.data)
}
