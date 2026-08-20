<template>
  <div class="page-container">
    <art-card :title="$t('multiRenew.title')" shadow="never">
      <!-- 操作栏 -->
      <div class="mb-4 flex justify-between">
        <el-button type="primary" @click="handleCreate">{{ $t('multiRenew.createTask') }}</el-button>
        <el-button @click="fetchData" :icon="Refresh">{{ $t('common.refresh') }}</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" :label="$t('multiRenew.taskId')" width="100" />
        <el-table-column prop="name" :label="$t('multiRenew.taskName')" min-width="150" />
        <el-table-column prop="service_count" :label="$t('multiRenew.serviceCount')" width="100" />
        <el-table-column prop="duration" :label="$t('multiRenew.duration')" width="100">
          <template #default="{ row }">{{ row.duration }}{{ $t('multiRenew.months') }}</template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" :label="$t('multiRenew.progress')" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.progress || 0" :status="row.status === 3 ? 'success' : undefined" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column prop="operator" :label="$t('common.operator')" width="120" />
        <el-table-column :label="$t('common.action')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="handleExecute(row)" :disabled="row.status !== 0">{{ $t('multiRenew.execute') }}</el-button>
            <el-button size="small" type="danger" @click="handleCancel(row)" :disabled="row.status !== 0 && row.status !== 1">{{ $t('common.cancel') }}</el-button>
            <el-button size="small" @click="handleViewLog(row)">{{ $t('multiRenew.logs') }}</el-button>
            <el-button size="small" type="info" @click="handleDelete(row)" :disabled="row.status === 1">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="mt-4 justify-end"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </art-card>

    <!-- 创建任务对话框 -->
    <el-dialog v-model="createDialogVisible" :title="$t('multiRenew.createTask')" width="600px" destroy-on-close>
      <el-form :model="createForm" :rules="createFormRules" ref="createFormRef" label-width="120px">
        <el-form-item :label="$t('multiRenew.taskName')" prop="name">
          <el-input v-model="createForm.name" :placeholder="$t('multiRenew.inputTaskName')" />
        </el-form-item>
        <el-form-item :label="$t('multiRenew.duration')" prop="duration">
          <el-select v-model="createForm.duration" :placeholder="$t('multiRenew.selectDuration')" style="width: 100%">
            <el-option :label="$t('multiRenew.oneMonth')" :value="1" />
            <el-option :label="$t('multiRenew.threeMonths')" :value="3" />
            <el-option :label="$t('multiRenew.sixMonths')" :value="6" />
            <el-option :label="$t('multiRenew.twelveMonths')" :value="12" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('multiRenew.filterType')" prop="filter_type">
          <el-radio-group v-model="createForm.filter_type">
            <el-radio value="ids">{{ $t('multiRenew.filterByIds') }}</el-radio>
            <el-radio value="product">{{ $t('multiRenew.filterByProduct') }}</el-radio>
            <el-radio value="due_date">{{ $t('multiRenew.filterByDueDate') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'ids'" :label="$t('multiRenew.serviceIdList')" prop="service_ids">
          <el-input v-model="createForm.service_ids" type="textarea" :rows="4" :placeholder="$t('multiRenew.inputServiceIds')" />
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'product'" :label="$t('multiRenew.selectProduct')" prop="product_id">
          <el-select v-model="createForm.product_id" :placeholder="$t('multiRenew.selectProductPlaceholder')" style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'due_date'" :label="$t('multiRenew.dueDateRange')" prop="due_date_range">
          <el-date-picker v-model="createForm.due_date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('common.remark')">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" :placeholder="$t('multiRenew.remarkPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleCreateSubmit" :loading="createSubmitLoading">{{ $t('multiRenew.createTask') }}</el-button>
      </template>
    </el-dialog>

    <!-- 日志抽屉 -->
    <el-drawer v-model="logDrawerVisible" :title="$t('multiRenew.taskLogs')" size="600px" destroy-on-close>
      <el-timeline>
        <el-timeline-item
          v-for="log in taskLogs"
          :key="log.id"
          :timestamp="log.created_at"
          :type="log.type === 'error' ? 'danger' : log.type === 'success' ? 'success' : 'primary'"
          placement="top"
        >
          <p>{{ log.content }}</p>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="taskLogs.length === 0" :description="$t('multiRenew.noLogs')" />
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { $t } from '@/locales'
import request from '@/utils/http'

defineOptions({ name: 'MultiRenew' })

const loading = ref(false)
const tableData = ref([])
const productList = ref<any[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 创建对话框
const createDialogVisible = ref(false)
const createSubmitLoading = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = reactive({
  name: '',
  duration: 1 as number,
  filter_type: 'ids' as 'ids' | 'product' | 'due_date',
  service_ids: '',
  product_id: undefined as number | undefined,
  due_date_range: null as any,
  remark: ''
})
const createFormRules: FormRules = {
  name: [{ required: true, message: $t('multiRenew.inputTaskName'), trigger: 'blur' }],
  duration: [{ required: true, message: $t('multiRenew.selectDuration'), trigger: 'change' }],
  filter_type: [{ required: true, message: $t('multiRenew.selectFilterType'), trigger: 'change' }]
}

// 日志抽屉
const logDrawerVisible = ref(false)
const taskLogs = ref<any[]>([])

const getStatusType = (status: number) => {
  const map: Record<number, any> = { 0: 'info', 1: 'warning', 2: 'danger', 3: 'success', 4: 'info' }
  return map[status] || 'info'
}

const statusTextMap: Record<number, () => string> = {
  0: () => $t('multiRenew.pending'),
  1: () => $t('multiRenew.processing'),
  2: () => $t('common.cancelled'),
  3: () => $t('common.completed'),
  4: () => $t('common.failed')
}

const getStatusText = (status: number) => statusTextMap[status]?.() || $t('common.unknown')

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/multi-renew',
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    tableData.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res = await request.get({ url: '/api/admin/products' })
    productList.value = res || []
  } catch (error) {
    console.error(error)
  }
}

const handleCreate = () => {
  createForm.name = ''
  createForm.duration = 1
  createForm.filter_type = 'ids'
  createForm.service_ids = ''
  createForm.product_id = undefined
  createForm.due_date_range = null
  createForm.remark = ''
  createDialogVisible.value = true
  fetchProducts()
}

const handleCreateSubmit = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    createSubmitLoading.value = true
    try {
      await request.post({ url: '/api/admin/multi-renew', params: createForm })
      ElMessage.success($t('multiRenew.taskCreated'))
      createDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('multiRenew.createFailed'))
    } finally {
      createSubmitLoading.value = false
    }
  })
}

const handleExecute = async (row: any) => {
  await ElMessageBox.confirm($t('multiRenew.confirmExecute'), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/multi-renew/${row.id}/execute` })
    ElMessage.success($t('multiRenew.taskExecuting'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleCancel = async (row: any) => {
  await ElMessageBox.confirm($t('multiRenew.confirmCancel'), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/multi-renew/${row.id}/cancel` })
    ElMessage.success($t('multiRenew.taskCancelled'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleViewLog = async (row: any) => {
  logDrawerVisible.value = true
  try {
    const res = await request.get({ url: `/api/admin/multi-renew/${row.id}/logs` })
    taskLogs.value = res || []
  } catch (error) {
    taskLogs.value = []
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('multiRenew.confirmDelete'), $t('common.warning'), { type: 'error' })
  try {
    await request.del({ url: `/api/admin/multi-renew/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
