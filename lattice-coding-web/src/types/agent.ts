export interface Agent {
  id: number
  name: string
  description: string
  model: string
  provider: string
  systemPrompt: string
  tools: string[]
  maxSteps: number
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface AgentForm {
  name: string
  description: string
  model: string
  provider: string
  systemPrompt: string
  tools: string[]
  maxSteps: number
  enabled: boolean
}
