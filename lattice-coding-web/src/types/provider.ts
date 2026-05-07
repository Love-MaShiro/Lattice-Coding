export interface Provider {
  id: number
  name: string
  provider_type: string
  base_url: string
  auth_type: string
  config: string
  enabled: boolean
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
