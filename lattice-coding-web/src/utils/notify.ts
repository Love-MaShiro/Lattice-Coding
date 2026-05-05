import { ElMessage, MessageOptions } from 'element-plus'

const DEFAULT_DURATION = 3000

const defaultOptions: Partial<MessageOptions> = {
  duration: DEFAULT_DURATION
}

export function notifySuccess(message: string, options?: Partial<MessageOptions>) {
  ElMessage.success({
    message,
    ...defaultOptions,
    ...options
  })
}

export function notifyError(message: string, options?: Partial<MessageOptions>) {
  ElMessage.error({
    message,
    ...defaultOptions,
    ...options
  })
}

export function notifyWarning(message: string, options?: Partial<MessageOptions>) {
  ElMessage.warning({
    message,
    ...defaultOptions,
    ...options
  })
}

export function notifyInfo(message: string, options?: Partial<MessageOptions>) {
  ElMessage.info({
    message,
    ...defaultOptions,
    ...options
  })
}
