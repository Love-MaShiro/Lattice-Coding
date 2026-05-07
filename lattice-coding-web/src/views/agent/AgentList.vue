<template>
  <PageContainer class="agent-page" title="Agent 管理" description="配置 Agent 的模型、提示词、上下文参数和工具权限">
    <template #actions>
      <el-button @click="refreshTable">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
      <el-button type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon>
        新增 Agent
      </el-button>
    </template>

    <LatticeTable ref="tableRef" :columns="columns" :api="loadAgents">
      <template #model="{ row }">
        <span class="strong-cell">{{ modelNameMap[row.model_config_id] || '-' }}</span>
      </template>
      <template #toolCount="{ row }">
        <span>{{ row.tool_count ?? 0 }}</span>
      </template>
      <template #temperature="{ row }">
        <span>{{ row.temperature ?? 0.7 }}</span>
      </template>
      <template #enabled="{ row }">
        <el-tag :type="row.enabled ? 'success' : 'info'" effect="plain">
          {{ row.enabled ? '启用' : '停用' }}
        </el-tag>
      </template>
      <template #createdAt="{ row }">
        <span>{{ formatTime(row.created_at) }}</span>
      </template>
      <template #actions="{ row }">
        <div class="row-actions">
          <el-button text type="primary" @click="openEdit(row)">
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button text :type="row.enabled ? 'warning' : 'success'" @click="toggleEnabled(row)">
            {{ row.enabled ? '停用' : '启用' }}
          </el-button>
          <el-button text type="danger" @click="removeAgent(row)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </template>
    </LatticeTable>

    <LatticeFormDialog
      ref="dialogRef"
      v-model="dialogVisible"
      title="Agent"
      width="760px"
      @submit="submitAgent"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="130px" class="agent-form">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入 Agent 名称" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="说明 Agent 的用途" />
        </el-form-item>

        <el-form-item label="模型" prop="model_config_id">
          <el-select
            v-model="form.model_config_id"
            placeholder="选择模型"
            filterable
            class="full-width"
            :loading="modelLoading"
          >
            <el-option-group
              v-for="group in groupedModels"
              :key="group.providerId"
              :label="group.providerName"
            >
              <el-option
                v-for="model in group.models"
                :key="model.id"
                :label="`${model.name} / ${model.model}`"
                :value="model.id"
              />
            </el-option-group>
          </el-select>
        </el-form-item>

        <el-form-item label="System Prompt">
          <el-input
            v-model="form.system_prompt"
            type="textarea"
            :rows="6"
            placeholder="输入 Agent 的系统提示词"
          />
        </el-form-item>

        <el-form-item label="temperature">
          <div class="slider-row">
            <el-slider v-model="form.temperature" :min="0" :max="1" :step="0.1" />
            <span class="slider-value">{{ form.temperature.toFixed(1) }}</span>
          </div>
        </el-form-item>

        <el-form-item label="top_p">
          <div class="slider-row">
            <el-slider v-model="form.top_p" :min="0" :max="1" :step="0.1" />
            <span class="slider-value">{{ form.top_p.toFixed(1) }}</span>
          </div>
        </el-form-item>

        <el-form-item label="max_tokens">
          <el-input-number v-model="form.max_tokens" :min="1" :max="200000" :step="512" />
        </el-form-item>

        <el-form-item label="max_context_turns">
          <el-input-number v-model="form.max_context_turns" :min="1" :max="100" />
        </el-form-item>

        <el-form-item label="max_steps">
          <el-input-number v-model="form.max_steps" :min="1" :max="1000" />
        </el-form-item>

        <el-form-item label="状态">
          <el-switch v-model="form.enabled" active-text="启用" inactive-text="停用" />
        </el-form-item>
      </el-form>
    </LatticeFormDialog>
  </PageContainer>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import dayjs from 'dayjs'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessageBox } from 'element-plus'
import { Delete, Edit, Plus, Refresh } from '@element-plus/icons-vue'

import { agentApi } from '@/api/agent'
import { providerApi } from '@/api/provider'
import LatticeFormDialog from '@/components/LatticeFormDialog.vue'
import LatticeTable, { type TableColumn } from '@/components/LatticeTable.vue'
import PageContainer from '@/components/PageContainer.vue'
import type { Agent, AgentForm } from '@/types/agent'
import type { ModelConfig, Provider } from '@/types/provider'
import { notifySuccess } from '@/utils/notify'

const columns: TableColumn[] = [
  { label: '名称', prop: 'name', minWidth: 180 },
  { label: '关联模型名', slot: 'model', minWidth: 260 },
  { label: '工具数量', slot: 'toolCount', minWidth: 110, align: 'center' },
  { label: 'temperature', slot: 'temperature', minWidth: 130, align: 'center' },
  { label: 'enabled', slot: 'enabled', minWidth: 110, align: 'center' },
  { label: '创建时间', slot: 'createdAt', minWidth: 190 },
  { label: '操作', slot: 'actions', width: 180, fixed: 'right', align: 'center' }
]

const tableRef = ref<{ refresh: () => void }>()
const dialogRef = ref<{
  open: (data?: Record<string, any>) => void
  setSubmitLoading: (loading: boolean) => void
}>()
const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const editingId = ref<number>()
const modelLoading = ref(false)
const providers = ref<Provider[]>([])
const modelConfigs = ref<ModelConfig[]>([])

const form = reactive<AgentForm>({
  name: '',
  description: '',
  model_config_id: undefined,
  system_prompt: '',
  temperature: 0.7,
  top_p: 1.0,
  max_tokens: 4096,
  max_context_turns: 10,
  max_steps: 20,
  enabled: true,
  tool_names: []
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入 Agent 名称', trigger: 'blur' }],
  model_config_id: [{ required: true, message: '请选择模型', trigger: 'change' }]
}

const providerNameMap = computed(() => {
  return providers.value.reduce<Record<number, string>>((acc, provider) => {
    acc[provider.id] = provider.name
    return acc
  }, {})
})

const modelNameMap = computed(() => {
  return modelConfigs.value.reduce<Record<number, string>>((acc, model) => {
    acc[model.id] = `${model.name} / ${model.model}`
    return acc
  }, {})
})

const groupedModels = computed(() => {
  const groups = new Map<number, { providerId: number; providerName: string; models: ModelConfig[] }>()
  modelConfigs.value
    .filter((model) => model.enabled)
    .forEach((model) => {
      if (!groups.has(model.provider_id)) {
        groups.set(model.provider_id, {
          providerId: model.provider_id,
          providerName: providerNameMap.value[model.provider_id] || `Provider #${model.provider_id}`,
          models: []
        })
      }
      groups.get(model.provider_id)?.models.push(model)
    })
  return Array.from(groups.values())
})

const normalizePageData = <T,>(result: any): T[] => {
  return result?.data || result?.list || result?.items || []
}

const loadModels = async () => {
  modelLoading.value = true
  try {
    const [providerPage, modelPage] = await Promise.all([
      providerApi.list({ page: 1, size: 100 }),
      providerApi.listModelConfigs({ page: 1, size: 100 } as any)
    ])
    providers.value = normalizePageData<Provider>(providerPage)
    modelConfigs.value = normalizePageData<ModelConfig>(modelPage)
  } finally {
    modelLoading.value = false
  }
}

const loadAgents = async (page: number, size: number) => {
  await loadModels()
  return agentApi.list({ page, page_size: size })
}

const resetForm = () => {
  editingId.value = undefined
  form.name = ''
  form.description = ''
  form.model_config_id = undefined
  form.system_prompt = ''
  form.temperature = 0.7
  form.top_p = 1.0
  form.max_tokens = 4096
  form.max_context_turns = 10
  form.max_steps = 20
  form.enabled = true
  form.tool_names = []
}

const openCreate = async () => {
  resetForm()
  await loadModels()
  dialogVisible.value = true
  dialogRef.value?.open()
}

const openEdit = async (agent: Agent) => {
  await loadModels()
  editingId.value = agent.id
  form.name = agent.name
  form.description = agent.description
  form.model_config_id = agent.model_config_id
  form.system_prompt = agent.system_prompt
  form.temperature = agent.temperature ?? 0.7
  form.top_p = agent.top_p ?? 1.0
  form.max_tokens = agent.max_tokens ?? 4096
  form.max_context_turns = agent.max_context_turns ?? 10
  form.max_steps = agent.max_steps ?? 20
  form.enabled = agent.enabled
  form.tool_names = []
  dialogVisible.value = true
  dialogRef.value?.open(agent as unknown as Record<string, any>)
}

const submitAgent = async () => {
  await formRef.value?.validate()
  dialogRef.value?.setSubmitLoading(true)
  try {
    if (editingId.value) {
      await agentApi.update(editingId.value, form)
      notifySuccess('Agent 已更新')
    } else {
      await agentApi.create(form)
      notifySuccess('Agent 已创建')
    }
    dialogVisible.value = false
    refreshTable()
  } finally {
    dialogRef.value?.setSubmitLoading(false)
  }
}

const toggleEnabled = async (agent: Agent) => {
  if (agent.enabled) {
    await agentApi.disable(agent.id)
    notifySuccess('Agent 已停用')
  } else {
    await agentApi.enable(agent.id)
    notifySuccess('Agent 已启用')
  }
  refreshTable()
}

const removeAgent = async (agent: Agent) => {
  await ElMessageBox.confirm(`确认删除 Agent「${agent.name}」？`, '删除 Agent', {
    type: 'warning',
    confirmButtonText: '删除',
    cancelButtonText: '取消'
  })
  await agentApi.delete(agent.id)
  notifySuccess('Agent 已删除')
  refreshTable()
}

const refreshTable = () => {
  tableRef.value?.refresh()
}

const formatTime = (value?: string) => {
  return value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '-'
}

onMounted(() => {
  loadModels()
})
</script>

<style scoped>
:deep(.page-content) {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 172px);
}

:deep(.lattice-table) {
  flex: 1;
}

.strong-cell {
  color: var(--color-text-primary);
  font-weight: 500;
}

.row-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.agent-form {
  padding-right: 8px;
}

.full-width {
  width: 100%;
}

.slider-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 44px;
  align-items: center;
  gap: 16px;
  width: 100%;
}

.slider-value {
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
  text-align: right;
}
</style>
