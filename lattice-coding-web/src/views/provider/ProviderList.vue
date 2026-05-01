<template>
  <div class="provider-list">
    <div v-if="status === 'loading'" class="status loading">正在检查后端连接...</div>
    <div v-else-if="status === 'success'" class="status success">后端已连接：Lattice-coding is running</div>
    <div v-else class="status error">后端未连接</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getHealth } from '@/api/health'

const status = ref<'loading' | 'success' | 'error'>('loading')

onMounted(async () => {
  try {
    await getHealth()
    status.value = 'success'
  } catch {
    status.value = 'error'
  }
})
</script>

<style scoped>
.status {
  font-size: 16px;
  padding: 10px;
}
.loading {
  color: #909399;
}
.success {
  color: #67c23a;
}
.error {
  color: #f56c6c;
}
</style>
