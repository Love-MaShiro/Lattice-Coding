export interface Provider {
  id: number
  name: string
  type: string
  baseURL: string
  apiKey: string
  models: string[]
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface ProviderForm {
  name: string
  type: string
  baseURL: string
  apiKey: string
  models: string[]
  enabled: boolean
}
