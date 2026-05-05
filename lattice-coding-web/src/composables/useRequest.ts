import { ref, shallowRef } from 'vue'

export interface UseRequestOptions {
  manual?: boolean
}

export function useRequest<T = any>(
  apiFn: () => Promise<T>,
  options: UseRequestOptions = {}
) {
  const { manual = false } = options

  const data = shallowRef<T | undefined>(undefined)
  const loading = ref(false)
  const error = shallowRef<Error | undefined>(undefined)

  const execute = async (): Promise<T | undefined> => {
    loading.value = true
    error.value = undefined

    try {
      const result = await apiFn()
      data.value = result
      return result
    } catch (err) {
      error.value = err as Error
      return undefined
    } finally {
      loading.value = false
    }
  }

  if (!manual) {
    execute()
  }

  return {
    data,
    loading,
    error,
    execute
  }
}
