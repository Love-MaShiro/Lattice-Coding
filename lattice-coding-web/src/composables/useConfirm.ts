import { ElMessageBox } from 'element-plus'
import { notifySuccess } from '@/utils/notify'

export interface ConfirmOptions {
  title?: string
  confirmText?: string
  cancelText?: string
  successMessage?: string
}

export async function confirmDelete<T = any>(
  apiFn: () => Promise<T>,
  options: ConfirmOptions = {}
): Promise<T | undefined> {
  const {
    title = '确认删除',
    confirmText = '删除',
    cancelText = '取消',
    successMessage = '删除成功'
  } = options

  try {
    await ElMessageBox.confirm(
      '删除后无法恢复，确定要删除吗？',
      title,
      {
        confirmButtonText: confirmText,
        cancelButtonText: cancelText,
        type: 'warning',
        draggable: true
      }
    )

    const result = await apiFn()
    notifySuccess(successMessage)
    return result
  } catch {
    // User cancelled or API error - do nothing
    return undefined
  }
}
