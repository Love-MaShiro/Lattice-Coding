<template>
  <PageContainer class="run-page" title="Run">
    <template #actions>
      <el-button :icon="Refresh" circle @click="loadRuns" />
    </template>

    <div class="run-layout">
      <section class="run-list">
        <el-table :data="runs" height="100%" v-loading="loading" highlight-current-row @row-click="selectRun">
          <el-table-column prop="id" label="Run ID" min-width="180" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="statusType(row.status)" effect="plain">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="Latency" width="110">
            <template #default="{ row }">{{ formatLatency(row) }}</template>
          </el-table-column>
          <el-table-column label="创建时间" width="150">
            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
          </el-table-column>
        </el-table>
      </section>

      <aside class="run-detail" v-loading="detailLoading">
        <template v-if="activeRun">
          <div class="detail-header">
            <div>
              <h3>{{ activeRun.id }}</h3>
              <el-tag :type="statusType(activeRun.status)" effect="plain">{{ activeRun.status }}</el-tag>
            </div>
          </div>

          <div class="metric-grid">
            <div class="metric-card">
              <span>Token</span>
              <strong>{{ activeRun.token_count ?? '待接入' }}</strong>
            </div>
            <div class="metric-card">
              <span>Latency</span>
              <strong>{{ formatLatency(activeRun) }}</strong>
            </div>
            <div class="metric-card">
              <span>Cost</span>
              <strong>{{ activeRun.cost != null ? `$${activeRun.cost.toFixed(6)}` : '待接入' }}</strong>
            </div>
          </div>

          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="Agent">{{ activeRun.agent_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="Session">{{ activeRun.session_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="Workflow">{{ activeRun.workflow_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="Started">{{ formatTime(activeRun.started_at) }}</el-descriptions-item>
            <el-descriptions-item label="Completed">{{ formatTime(activeRun.completed_at) }}</el-descriptions-item>
          </el-descriptions>

          <div class="detail-block">
            <h4>Input</h4>
            <pre>{{ activeRun.input || '-' }}</pre>
          </div>
          <div class="detail-block">
            <h4>Output</h4>
            <pre>{{ activeRun.output || '-' }}</pre>
          </div>
          <div v-if="activeRun.error" class="detail-block error">
            <h4>Error</h4>
            <pre>{{ activeRun.error }}</pre>
          </div>

          <div class="detail-block">
            <h4>Tool Invocations</h4>
            <el-table :data="toolInvocations" size="small" empty-text="暂无工具调用">
              <el-table-column prop="tool_name" label="Tool" min-width="140" show-overflow-tooltip />
              <el-table-column prop="status" label="状态" width="90" />
              <el-table-column prop="latency_ms" label="Latency" width="100" />
            </el-table>
          </div>
        </template>
        <el-empty v-else description="选择一个 Run 查看详情" />
      </aside>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { Refresh } from '@element-plus/icons-vue'

import { runApi } from '@/api/run'
import PageContainer from '@/components/PageContainer.vue'
import type { Run, ToolInvocation } from '@/types/run'
import { notifyError } from '@/utils/notify'

const loading = ref(false)
const detailLoading = ref(false)
const runs = ref<Run[]>([])
const activeRunId = ref<string>()
const activeRun = computed(() => runs.value.find((run) => run.id === activeRunId.value))
const toolInvocations = ref<ToolInvocation[]>([])

async function loadRuns() {
  loading.value = true
  try {
    const result = await runApi.list({ page: 1, size: 50 })
    runs.value = result.data || []
    if (!activeRunId.value && runs.value.length > 0) {
      await selectRun(runs.value[0])
    }
  } catch (err) {
    notifyError(err instanceof Error ? err.message : '加载 Run 失败')
  } finally {
    loading.value = false
  }
}

async function selectRun(run: Run) {
  activeRunId.value = run.id
  detailLoading.value = true
  try {
    const [detail, invocations] = await Promise.all([runApi.get(run.id), runApi.listToolInvocations(run.id)])
    runs.value = runs.value.map((item) => (item.id === detail.id ? detail : item))
    toolInvocations.value = invocations
  } catch (err) {
    notifyError(err instanceof Error ? err.message : '加载 Run 详情失败')
  } finally {
    detailLoading.value = false
  }
}

function statusType(status: string) {
  if (status === 'completed') return 'success'
  if (status === 'failed') return 'danger'
  if (status === 'running') return 'warning'
  return 'info'
}

function formatLatency(run: Run) {
  if (run.latency_ms != null) return `${run.latency_ms} ms`
  if (!run.started_at || !run.completed_at) return '待接入'
  const ms = dayjs(run.completed_at).diff(dayjs(run.started_at))
  return ms >= 0 ? `${ms} ms` : '待接入'
}

function formatTime(value?: string) {
  return value ? dayjs(value).format('MM-DD HH:mm:ss') : '-'
}

onMounted(loadRuns)
</script>

<style scoped>
:deep(.page-content) {
  min-height: calc(100vh - 172px);
  padding: 0;
}

.run-layout {
  min-height: calc(100vh - 172px);
  display: grid;
  grid-template-columns: minmax(0, 1fr) 420px;
}

.run-list {
  min-width: 0;
  padding: 16px;
  border-right: 1px solid var(--color-border);
}

.run-detail {
  min-width: 0;
  padding: 16px;
  overflow: auto;
  background: #f8fafc;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.detail-header h3 {
  margin: 0 0 8px;
  font-size: 16px;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 16px;
}

.metric-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: #fff;
}

.metric-card span {
  color: var(--color-text-secondary);
  font-size: 12px;
}

.metric-card strong {
  color: var(--color-text-primary);
  font-size: 16px;
}

.detail-block {
  margin-top: 16px;
}

.detail-block h4 {
  margin: 0 0 8px;
  font-size: 14px;
}

.detail-block pre {
  min-height: 48px;
  max-height: 220px;
  margin: 0;
  padding: 12px;
  overflow: auto;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: #fff;
  white-space: pre-wrap;
  word-break: break-word;
}

.detail-block.error pre {
  border-color: #fecaca;
  color: #991b1b;
  background: #fef2f2;
}
</style>
