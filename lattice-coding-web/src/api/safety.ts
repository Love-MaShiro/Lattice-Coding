import { post } from './request'

export interface SafetyCheckResponse {
  allowed: boolean
  reason?: string
}

export const safetyApi = {
  checkPath: (path: string) =>
    post<SafetyCheckResponse>('/safety/check/path', { path }),

  checkCommand: (command: string) =>
    post<SafetyCheckResponse>('/safety/check/command', { command }),

  checkPermission: (permission: string) =>
    post<SafetyCheckResponse>('/safety/check/permission', { permission })
}
