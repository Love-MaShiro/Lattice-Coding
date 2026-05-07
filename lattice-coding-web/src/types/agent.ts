export interface AgentTool {
  id: number
  tool_id: number
  tool_type: string
  created_at: string
}

export interface Agent {
  id: number
  name: string
  description: string
  agent_type: string
  model_config_id: number
  system_prompt: string
  temperature: number
  top_p: number
  max_tokens: number
  max_context_turns: number
  max_steps: number
  enabled: boolean
  tool_count: number
  created_at: string
  updated_at: string
}

export interface AgentDetail extends Agent {
  tools: AgentTool[]
}

export interface AgentForm {
  name: string
  description: string
  model_config_id: number | undefined
  system_prompt: string
  temperature: number
  top_p: number
  max_tokens: number
  max_context_turns: number
  max_steps: number
  enabled: boolean
  tool_names: string[]
}

export const toolTypeMap: Record<string, string> = {
  shell: 'local',
  file: 'local',
  git: 'local',
  mcp: 'mcp',
  knowledge: 'knowledge'
}
