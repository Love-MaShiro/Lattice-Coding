import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const collapsed = ref(false)
  const title = ref('控制台')

  const toggleCollapse = () => {
    collapsed.value = !collapsed.value
  }

  return {
    collapsed,
    title,
    toggleCollapse
  }
})
