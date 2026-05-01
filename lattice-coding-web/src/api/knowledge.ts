import { get, post, put, del } from './request'
import type { KnowledgeDocument, KnowledgeForm } from '@/types/knowledge'
import type { PageResult, PageQuery } from '@/types/api'

export const knowledgeApi = {
  list: (params?: PageQuery) =>
    get<PageResult<KnowledgeDocument>>('/knowledge/documents', { params }),

  get: (id: number) => get<KnowledgeDocument>(`/knowledge/documents/${id}`),

  create: (data: KnowledgeForm) =>
    post<KnowledgeDocument>('/knowledge/documents', data),

  update: (id: number, data: KnowledgeForm) =>
    put<KnowledgeDocument>(`/knowledge/documents/${id}`, data),

  delete: (id: number) => del<void>(`/knowledge/documents/${id}`),

  search: (query: string) =>
    post<KnowledgeDocument[]>('/knowledge/search', { query })
}
