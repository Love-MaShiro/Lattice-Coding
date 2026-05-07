<template>
  <div class="model-config-panel">
    <div class="panel-header">
      <span class="panel-title">模型配置管理</span>
      <el-button type="primary" size="small" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增模型配置
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="configs"
      border
      stripe
      size="small"
      style="width: 100%"
    >
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="model" label="模型名" min-width="140" />
      <el-table-column prop="model_type" label="模型类型" width="100">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.model_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_default" label="默认" width="70" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" size="small" type="success">默认</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="enabled" label="启用" width="70" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.enabled"
            size="small"
            @change="(val: boolean) => handleToggleEnabled(row, val)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" width="170">
        <template #default="{ row }">
          {{ formatTime(row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button
            v-if="!row.is_default"
            link
            type="warning"
            size="small"
            @click="handleSetDefault(row)"
          >
            设为默认
          </el-button>
          <el-button
            link
            type="primary"
            size="small"
            :loading="testingId === row.id"
            @click="handleTest(row)"
          >
            测试
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty description="暂无模型配置" :image-size="80" />
      </template>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑模型配置' : '新增模型配置'"
      width="600px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入配置名称" />
        </el-form-item>
        <el-form-item label="模型名" prop="model">
          <el-input v-model="form.model" placeholder="例如 gpt-4o-mini、deepseek-chat" />
        </el-form-item>
        <el-form-item label="模型类型" prop="model_type">
          <el-select v-model="form.model_type" placeholder="请选择模型类型" style="width: 100%">
            <el-option label="chat" value="chat" />
            <el-option label="embedding" value="embedding" />
            <el-option label="rerank" value="rerank" />
          </el-select>
        </el-form-item>
        <el-form-item label="参数" prop="params">
          <el-input
            v-model="form.params"
            type="textarea"
            :rows="6"
            placeholder='{"temperature": 0.7, "top_p": 0.9, "max_tokens": 2048, "timeout": 60}'
          />
        </el-form-item>
        <el-form-item label="能力" prop="capabilities">
          <el-input
            v-model="form.capabilities"
            type="textarea"
            :rows="4"
            placeholder='{"stream": true, "tool_call": false, "vision": false, "json_mode": false}'
          />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.is_default" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { providerApi } from '@/api/provider'
import type { ModelConfig, ModelConfigForm } from '@/types/provider'
import { notifySuccess, notifyError } from '@/utils/notify'

const props = defineProps<{
  providerId: number
}>()

const loading = ref(false)
const configs = ref<ModelConfig[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const submitting = ref(false)
const testingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const defaultParams = JSON.stringify(
  {
    temperature: 0.7,
    top_p: 0.9,
    max_tokens: 2048,
    timeout: 60
  },
  null,
  2
)

const defaultCapabilities = JSON.stringify(
  {
    stream: true,
    tool_call: false,
    vision: false,
    json_mode: false
  },
  null,
  2
)

const form = reactive<ModelConfigForm & { is_default: boolean }>({
  provider_id: props.providerId,
  name: '',
  model: '',
  model_type: 'chat',
  params: defaultParams,
  capabilities: defaultCapabilities,
  is_default: false,
  enabled: true
})

const validateJson = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value || value.trim() === '') {
    callback()
    return
  }
  try {
    JSON.parse(value)
    callback()
  } catch {
    callback(new Error('请输入合法的 JSON 格式'))
  }
}

const formRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  model: [{ required: true, message: '请输入模型名', trigger: 'blur' }],
  model_type: [{ required: true, message: '请选择模型类型', trigger: 'change' }],
  params: [{ validator: validateJson, trigger: 'blur' }],
  capabilities: [{ validator: validateJson, trigger: 'blur' }]
}

function formatTime(time: string) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

async function loadConfigs() {
  loading.value = true
  try {
    const result = await providerApi.listModelConfigs({
      provider_id: props.providerId,
      page: 1,
      size: 100
    })
    configs.value = result.data || []
  } catch {
    configs.value = []
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.provider_id = props.providerId
  form.name = ''
  form.model = ''
  form.model_type = 'chat'
  form.params = defaultParams
  form.capabilities = defaultCapabilities
  form.is_default = false
  form.enabled = true
  formRef.value?.resetFields()
}

function handleCreate() {
  isEdit.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: ModelConfig) {
  isEdit.value = true
  editingId.value = row.id
  form.provider_id = props.providerId
  form.name = row.name
  form.model = row.model
  form.model_type = row.model_type
  form.params = row.params || defaultParams
  form.capabilities = row.capabilities || defaultCapabilities
  form.is_default = row.is_default
  form.enabled = row.enabled
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const payload: ModelConfigForm = {
      provider_id: form.provider_id,
      name: form.name,
      model: form.model,
      model_type: form.model_type,
      params: form.params,
      capabilities: form.capabilities,
      is_default: form.is_default,
      enabled: form.enabled
    }

    if (isEdit.value && editingId.value) {
      await providerApi.updateModelConfig(editingId.value, payload)
      notifySuccess('模型配置更新成功')
    } else {
      await providerApi.createModelConfig(payload)
      notifySuccess('模型配置创建成功')
    }
    dialogVisible.value = false
    await loadConfigs()
  } catch {
    notifyError(isEdit.value ? '模型配置更新失败' : '模型配置创建失败')
  } finally {
    submitting.value = false
  }
}

async function handleToggleEnabled(row: ModelConfig, enabled: boolean) {
  try {
    if (enabled) {
      await providerApi.enableModelConfig(row.id)
    } else {
      await providerApi.disableModelConfig(row.id)
    }
    notifySuccess(enabled ? '已启用' : '已禁用')
    await loadConfigs()
  } catch {
    notifyError('操作失败')
  }
}

async function handleSetDefault(row: ModelConfig) {
  try {
    await ElMessageBox.confirm(`确定将 "${row.name}" 设为默认模型配置吗？`, '确认操作', {
      type: 'warning'
    })
    await providerApi.setDefaultModelConfig(row.id)
    notifySuccess('已设为默认')
    await loadConfigs()
  } catch {
    // cancelled
  }
}

async function handleTest(row: ModelConfig) {
  testingId.value = row.id
  try {
    const result = await providerApi.testModelConfig(row.id)
    if (result.success) {
      notifySuccess(`模型连接成功，耗时 ${result.latency_ms}ms`)
    } else {
      notifyError(result.error || '模型连接失败')
    }
  } catch {
    notifyError('模型测试失败')
  } finally {
    testingId.value = null
  }
}

async function handleDelete(row: ModelConfig) {
  try {
    await ElMessageBox.confirm(`确定要删除模型配置 "${row.name}" 吗？此操作不可恢复。`, '确认删除', {
      type: 'warning',
      confirmButtonText: '确定删除',
      cancelButtonText: '取消'
    })
    await providerApi.deleteModelConfig(row.id)
    notifySuccess('删除成功')
    await loadConfigs()
  } catch {
    // cancelled
  }
}

watch(
  () => props.providerId,
  () => {
    if (props.providerId) {
      loadConfigs()
    }
  },
  { immediate: true }
)

defineExpose({ refresh: loadConfigs })
</script>

<style scoped>
.model-config-panel {
  padding: 16px 0;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.panel-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}
</style>
