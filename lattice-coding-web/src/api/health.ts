import { get } from './request'

export interface HealthResponse {
  code: number | string
  message: string
  data: string
}

export function getHealth(): Promise<HealthResponse> {
  return get<HealthResponse>('/health')
}
