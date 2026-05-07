export interface Provider {
  id: number
  name: string
  provider_type: string
  base_url: string
  auth_type: string
  api_key_set: boolean
  config: string
  enabled: boolean
  health_status: string
  last_checked_at: string | null
  last_error: string
  created_at: string
  updated_at: string
}

export interface ProviderForm {
  name: string
  provider_type: string
  base_url: string
  auth_type: string
  api_key?: string
  auth_config?: string
  config?: string
  enabled: boolean
}

export interface ModelConfig {
  id: number
  provider_id: number
  provider_name?: string
  name: string
  model: string
  model_type: string
  params: string
  capabilities: string
  is_default: boolean
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface ModelConfigForm {
  provider_id: number
  name: string
  model: string
  model_type: string
  params?: string
  capabilities?: string
  is_default?: boolean
  enabled: boolean
}

export interface ProviderTestResult {
  success: boolean
  latency_ms: number
  error?: string
}

export interface ProviderHealthResult {
  provider_id: number
  model_config_id?: number
  status: string
  latency_ms: number
  error_code?: string
  error_message?: string
  checked_at: string
}

export interface SyncModelsResult {
  provider_id: number
  total: number
  created: number
  skipped: number
  failed: number
  message?: string
}

export interface ModelTestResult {
  success: boolean
  latency_ms: number
  error?: string
}
