export interface Run {
  id: number
  agentId: number
  agentName: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  input: string
  output: string
  error: string
  startedAt: string
  finishedAt: string
  createdAt: string
}

export interface RunForm {
  agentId: number
  input: string
}
