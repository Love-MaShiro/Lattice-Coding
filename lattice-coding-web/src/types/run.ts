export interface Run {
  id: string
  agent_id?: string
  session_id?: string
  workflow_id?: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled' | string
  input?: string
  output?: string
  error?: string
  started_at?: string
  completed_at?: string
  created_at?: string
  updated_at?: string
  token_count?: number
  latency_ms?: number
  cost?: number
}

export interface ToolInvocation {
  id: string
  run_id?: string
  node_id?: string
  tool_name: string
  input_json?: string
  output_json?: string
  is_error: boolean
  latency_ms: number
  status: string
  full_result_ref?: string
  started_at?: string
  completed_at?: string
  created_at?: string
  updated_at?: string
}

export interface RunForm {
  agent_id: number
  input: string
}
