<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    :width="width"
    @closed="handleClosed"
  >
    <slot></slot>
    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
        {{ mode === 'create' ? '创建' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref, computed } from 'vue'
import type { FormRules } from 'element-plus'

export type DialogMode = 'create' | 'edit'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title: string
  width?: string | number
  rules?: FormRules
}>(), {
  width: '500px'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [formData: T, mode: DialogMode]
}>()

const formData = ref<T>({} as T)
const mode = ref<DialogMode>('create')
const submitLoading = ref(false)

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const dialogTitle = computed(() => {
  return mode.value === 'create' ? `新增 ${props.title}` : `编辑 ${props.title}`
})

const handleClosed = () => {
  formData.value = {} as T
}

const open = (data?: T) => {
  if (data) {
    mode.value = 'edit'
    formData.value = { ...data }
  } else {
    mode.value = 'create'
    formData.value = {} as T
  }
}

const handleSubmit = () => {
  emit('submit', formData.value, mode.value)
}

const handleCancel = () => {
  visible.value = false
}

const setSubmitLoading = (loading: boolean) => {
  submitLoading.value = loading
}

defineExpose({
  open,
  setSubmitLoading
})
</script>

<style scoped>
</style>
