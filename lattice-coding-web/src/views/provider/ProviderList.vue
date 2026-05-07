<template>
  <PageContainer title="模型供应商管理">
    <template #actions>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增 Provider
      </el-button>
      <el-button @click="refresh">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </template>

    <el-table
      v-loading="loading"
      :data="providers"
      border
      stripe
      style="width: 100%"
      row-key="id"
    >
      <el-table-column type="expand">
        <template #default="{ row }">
          <ModelConfigPanel
            :ref="(el: any) => setModelConfigRef(row.id, el)"
            :provider-id="row.id"
          />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="provider_type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag size="small">{{ row.provider_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="base_url" label="Base URL" min-width="200" show-overflow-tooltip />
      <el-table-column prop="auth_type" label="鉴权方式" width="110">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.auth_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="API Key" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.api_key_set ? 'success' : 'info'" size="small">
            {{ row.api_key_set ? '已配置' : '未配置' }}
          </el-tag>
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
      <el-table-column label="健康状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="healthTagType(row.health_status)" size="small">
            {{ healthLabel(row.health_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_checked_at" label="最近检测" width="170">
        <template #default="{ row }">
          {{ formatTime(row.last_checked_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="340" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button
            link
            type="primary"
            size="small"
            :loading="testingId === row.id"
            @click="handleTest(row)"
          >
            测试
          </el-button>
          <el-button
            link
            type="primary"
            size="small"
            :loading="syncingId === row.id"
            @click="handleSyncModels(row)"
          >
            同步模型
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty description="暂无 Provider，请点击新增" :image-size="80" />
      </template>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑 Provider' : '新增 Provider'"
      width="650px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="110px"
        label-position="right"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入 Provider 名称" />
        </el-form-item>
        <el-form-item label="类型" prop="provider_type">
          <el-select
            v-model="form.provider_type"
            placeholder="请选择 Provider 类型"
            style="width: 100%"
            @change="onProviderTypeChange"
          >
            <el-option label="openai" value="openai" />
            <el-option label="openai_compatible" value="openai_compatible" />
            <el-option label="deepseek" value="deepseek" />
            <el-option label="qwen" value="qwen" />
            <el-option label="ollama" value="ollama" />
            <el-option label="claude" value="claude" />
          </el-select>
        </el-form-item>
        <el-form-item label="Base URL" prop="base_url">
          <el-input v-model="form.base_url" placeholder="请输入 Base URL" />
        </el-form-item>
        <el-form-item label="鉴权方式" prop="auth_type">
          <el-select v-model="form.auth_type" placeholder="请选择鉴权方式" style="width: 100%">
            <el-option label="none" value="none" />
            <el-option label="bearer" value="bearer" />
            <el-option label="api_key" value="api_key" />
            <el-option label="custom_header" value="custom_header" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input
            v-model="form.api_key"
            type="password"
            show-password
            :placeholder="isEdit ? '已配置 API Key，留空表示不修改' : '请输入 API Key'"
          />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        <el-form-item label="Config" prop="config">
          <el-input
            v-model="form.config"
            type="textarea"
            :rows="3"
            placeholder='可选，JSON 格式，例如 {"timeout": 30}'
          />
        </el-form-item>
        <el-form-item label="Auth Config" prop="auth_config">
          <el-input
            v-model="form.auth_config"
            type="textarea"
            :rows="3"
            placeholder='可选，JSON 格式，例如 {"header_name": "X-API-Key"}'
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </PageContainer>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import PageContainer from '@/components/PageContainer.vue'
import ModelConfigPanel from './ModelConfigPanel.vue'
import { providerApi } from '@/api/provider'
import type { Provider, ProviderForm } from '@/types/provider'
import { notifySuccess, notifyError } from '@/utils/notify'

const DEFAULT_BASE_URLS: Record<string, string> = {
  openai: 'https://api.openai.com/v1',
  openai_compatible: '',
  deepseek: 'https://api.deepseek.com/v1',
  qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
  ollama: 'http://localhost:11434/v1',
  claude: 'https://api.anthropic.com/v1'
}

const loading = ref(false)
const providers = ref<Provider[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const submitting = ref(false)
const testingId = ref<number | null>(null)
const syncingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const modelConfigRefs = ref<Record<number, any>>({})

const form = reactive<ProviderForm & { api_key: string; auth_config: string; config: string }>({
  name: '',
  provider_type: '',
  base_url: '',
  auth_type: 'api_key',
  api_key: '',
  auth_config: '',
  config: '',
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
  provider_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  auth_type: [{ required: true, message: '请选择鉴权方式', trigger: 'change' }],
  config: [{ validator: validateJson, trigger: 'blur' }],
  auth_config: [{ validator: validateJson, trigger: 'blur' }]
}

function healthTagType(status: string): 'success' | 'warning' | 'danger' | 'info' {
  switch (status) {
    case 'healthy':
      return 'success'
    case 'degraded':
      return 'warning'
    case 'unhealthy':
      return 'danger'
    default:
      return 'info'
  }
}

function healthLabel(status: string): string {
  switch (status) {
    case 'healthy':
      return '健康'
    case 'degraded':
      return '降级'
    case 'unhealthy':
      return '异常'
    default:
      return '未知'
  }
}

function formatTime(time: string | null) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

function setModelConfigRef(id: number, el: any) {
  if (el) {
    modelConfigRefs.value[id] = el
  }
}

function onProviderTypeChange(type: string) {
  if (!form.base_url && DEFAULT_BASE_URLS[type]) {
    form.base_url = DEFAULT_BASE_URLS[type]
  }
}

async function loadProviders() {
  loading.value = true
  try {
    const result = await providerApi.list({ page: 1, size: 100 })
    providers.value = result.data || []
  } catch {
    providers.value = []
  } finally {
    loading.value = false
  }
}

function refresh() {
  loadProviders()
}

function resetForm() {
  form.name = ''
  form.provider_type = ''
  form.base_url = ''
  form.auth_type = 'api_key'
  form.api_key = ''
  form.auth_config = ''
  form.config = ''
  form.enabled = true
  formRef.value?.resetFields()
}

function handleCreate() {
  isEdit.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: Provider) {
  isEdit.value = true
  editingId.value = row.id
  form.name = row.name
  form.provider_type = row.provider_type
  form.base_url = row.base_url
  form.auth_type = row.auth_type
  form.api_key = ''
  form.auth_config = ''
  form.config = row.config || ''
  form.enabled = row.enabled
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const payload: ProviderForm = {
      name: form.name,
      provider_type: form.provider_type,
      base_url: form.base_url,
      auth_type: form.auth_type,
      enabled: form.enabled,
      config: form.config || undefined,
      auth_config: form.auth_config || undefined
    }

    if (form.api_key) {
      payload.api_key = form.api_key
    }

    if (isEdit.value && editingId.value) {
      await providerApi.update(editingId.value, payload)
      notifySuccess('Provider 更新成功')
    } else {
      await providerApi.create(payload)
      notifySuccess('Provider 创建成功')
    }
    dialogVisible.value = false
    await loadProviders()
  } catch {
    notifyError(isEdit.value ? 'Provider 更新失败' : 'Provider 创建失败')
  } finally {
    submitting.value = false
  }
}

async function handleToggleEnabled(row: Provider, enabled: boolean) {
  try {
    if (enabled) {
      await providerApi.enable(row.id)
    } else {
      await providerApi.disable(row.id)
    }
    notifySuccess(enabled ? '已启用' : '已禁用')
    await loadProviders()
  } catch {
    notifyError('操作失败')
  }
}

async function handleTest(row: Provider) {
  testingId.value = row.id
  try {
    const result = await providerApi.test(row.id)
    if (result.success) {
      notifySuccess(`连接成功，耗时 ${result.latency_ms}ms`)
    } else {
      notifyError(result.error || '连接失败')
    }
    await loadProviders()
  } catch {
    notifyError('测试失败')
  } finally {
    testingId.value = null
  }
}

async function handleSyncModels(row: Provider) {
  syncingId.value = row.id
  try {
    const result = await providerApi.syncModels(row.id)
    const parts: string[] = []
    if (result.total > 0) parts.push(`共 ${result.total} 个模型`)
    if (result.created > 0) parts.push(`新增 ${result.created}`)
    if (result.skipped > 0) parts.push(`跳过 ${result.skipped}`)
    if (result.failed > 0) parts.push(`失败 ${result.failed}`)
    notifySuccess(parts.length > 0 ? parts.join('，') : (result.message || '同步完成'))
    const mcRef = modelConfigRefs.value[row.id]
    if (mcRef?.refresh) {
      mcRef.refresh()
    }
  } catch {
    notifyError('模型同步失败')
  } finally {
    syncingId.value = null
  }
}

async function handleDelete(row: Provider) {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Provider "${row.name}" 吗？删除后关联的模型配置也将被删除，此操作不可恢复。`,
      '确认删除',
      {
        type: 'warning',
        confirmButtonText: '确定删除',
        cancelButtonText: '取消'
      }
    )
    await providerApi.delete(row.id)
    notifySuccess('删除成功')
    await loadProviders()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadProviders()
})
</script>
