<template>
  <PageContainer title="模型管理" description="管理 AI 模型提供商配置">
    <template #actions>
      <el-button type="primary" class="btn-primary-gradient">
        <el-icon><Plus /></el-icon>
        新增 Provider
      </el-button>
    </template>

    <div v-if="status === 'loading'" class="status loading">正在检查后端连接...</div>
    <div v-else-if="status === 'success'" class="status success">后端已连接：Lattice-coding is running</div>
    <div v-else class="status error">后端未连接</div>
  </PageContainer>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getHealth } from '@/api/health'
import PageContainer from '@/components/PageContainer.vue'
import { Plus } from '@element-plus/icons-vue'

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
  color: var(--color-text-secondary);
}
.success {
  color: var(--color-success-600);
}
.error {
  color: var(--color-error-600);
}

.btn-primary-gradient {
  background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
  border: none;
  color: #fff;
  font-weight: 500;
  transition: all var(--transition-normal);
}

.btn-primary-gradient:hover {
  background: linear-gradient(135deg, var(--color-primary-600), var(--color-primary-700));
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.btn-secondary-outline {
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  transition: all var(--transition-normal);
}

.btn-secondary-outline:hover {
  border-color: var(--color-primary-400);
  color: var(--color-primary-600);
}

.btn-danger {
  background: var(--color-error-600);
  border: none;
  color: #fff;
  transition: all var(--transition-normal);
}

.btn-danger:hover {
  background: var(--color-error-500);
}
</style>
