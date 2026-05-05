<template>
  <div class="lattice-table">
    <el-table
      v-loading="loading"
      :data="tableData"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column
        v-for="col in columns"
        :key="col.prop || col.slot"
        :label="col.label"
        :prop="col.prop"
        :width="col.width"
      >
        <template #default="{ row }" v-if="col.slot">
          <slot :name="col.slot" :row="row" />
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" description="暂无数据" />

    <div class="pagination-wrapper" v-if="showPagination && total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="size"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref, watch, onMounted } from 'vue'

export interface TableColumn {
  label: string
  prop?: string
  width?: string | number
  slot?: string
}

export interface PageResult {
  items: T[]
  total: number
  page: number
  size: number
}

export type LoadDataFn = (page: number, size: number) => Promise<PageResult>

const props = withDefaults(defineProps<{
  columns: TableColumn[]
  api: LoadDataFn
  showPagination?: boolean
}>(), {
  showPagination: true
})

const emit = defineEmits<{
  refresh: []
}>()

const loading = ref(false)
const tableData = ref<T[]>([])
const page = ref(1)
const size = ref(20)
const total = ref(0)

const loadData = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const result = await props.api(page.value, size.value)
    tableData.value = result.items
    total.value = result.total
  } catch (error) {
    console.error('Failed to load data:', error)
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handlePageChange = () => {
  loadData()
}

const handleSizeChange = () => {
  page.value = 1
  loadData()
}

const refresh = () => {
  loadData()
  emit('refresh')
}

defineExpose({
  refresh
})

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.lattice-table {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
}
</style>
