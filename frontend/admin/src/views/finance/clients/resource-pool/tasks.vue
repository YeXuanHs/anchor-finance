<template>
  <div class="resource-pool-tasks-page">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsResourcePool.taskType')">
          <el-select v-model="searchForm.task_type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsResourcePool.taskAllocate')" value="allocate" />
            <el-option :label="$t('clientsResourcePool.taskReclaim')" value="reclaim" />
            <el-option :label="$t('clientsResourcePool.taskMigrate')" value="migrate" />
            <el-option :label="$t('clientsResourcePool.taskExpand')" value="expand" />
            <el-option :label="$t('clientsResourcePool.taskHealthCheck')" value="health_check" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsResourcePool.taskStatus')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsResourcePool.statusPending')" value="pending" />
            <el-option :label="$t('clientsResourcePool.statusRunning')" value="running" />
            <el-option :label="$t('clientsResourcePool.statusCompleted')" value="completed" />
            <el-option :label="$t('clientsResourcePool.statusFailed')" value="failed" />
            <el-option :label="$t('clientsResourcePool.statusCancelled')" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('common.search') }}
          </el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="art-table-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsResourcePool.taskQueue') }}</span>
          <el-space>
            <el-tag type="warning">{{ $t('clientsResourcePool.statusPending') }}: {{ stats.pending || 0 }}</el-tag>
            <el-tag type="primary">{{ $t('clientsResourcePool.statusRunning') }}: {{ stats.running || 0 }}</el-tag>
            <el-tag type="danger">{{ $t('clientsResourcePool.statusFailed') }}: {{ stats.failed || 0 }}</el-tag>
          </el-space>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="task_no" :label="$t('clientsResourcePool.taskNo')" width="150" />
        <el-table-column prop="task_type" :label="$t('clientsResourcePool.taskType')" width="120">
          <template #default="{ row }">
            <el-tag :type="getTaskTypeTag(row.task_type)" size="small">
              {{ getTaskTypeText(row.task_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_resource" :label="$t('clientsResourcePool.targetResource')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" :label="$t('clientsResourcePool.progress')" width="120">
          <template #default="{ row }">
            <el-progress
              :percentage="row.progress || 0"
              :status="row.status === 'completed' ? 'success' : row.status === 'failed' ? 'exception' : undefined"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column prop="completed_at" :label="$t('clientsResourcePool.completedAt')" width="180">
          <template #default="{ row }">
            {{ row.completed_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="operator" :label="$t('common.operator')" width="100" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('common.detail') }}</el-button>
            <template v-if="row.status === 'failed'">
              <el-popconfirm :title="$t('clientsResourcePool.confirmRetry')" @confirm="handleRetry(row)">
                <template #reference>
                  <el-button type="warning" link>{{ $t('clientsResourcePool.retry') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
            <template v-if="row.status === 'pending' || row.status === 'running'">
              <el-popconfirm :title="$t('clientsResourcePool.confirmCancel')" @confirm="handleCancel(row)">
                <template #reference>
                  <el-button type="danger" link>{{ $t('clientsResourcePool.cancel') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="detailDialogVisible" :title="$t('clientsResourcePool.taskDetail')" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('clientsResourcePool.taskNo')">{{ detailData.task_no }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.taskType')">{{ getTaskTypeText(detailData.task_type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.targetResource')" :span="2">{{ detailData.target_resource }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.status')">
          <el-tag :type="getStatusTag(detailData.status)">{{ getStatusText(detailData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.progress')">{{ detailData.progress || 0 }}%</el-descriptions-item>
        <el-descriptions-item :label="$t('common.createdAt')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.completedAt')">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.operator')">{{ detailData.operator || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.duration')">{{ detailData.duration || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.params')" :span="2">
          <pre class="detail-pre">{{ detailData.params ? JSON.stringify(detailData.params, null, 2) : '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.result')" :span="2">
          <pre class="detail-pre">{{ detailData.result ? JSON.stringify(detailData.result, null, 2) : '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('clientsResourcePool.errorMsg')" :span="2" v-if="detailData.error_message">
          <span class="text-danger">{{ detailData.error_message }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <el-divider v-if="detailData.logs && detailData.logs.length" />
      <div v-if="detailData.logs && detailData.logs.length" class="task-logs">
        <h4>{{ $t('clientsResourcePool.executionLogs') }}</h4>
        <el-timeline>
          <el-timeline-item
            v-for="(log, index) in detailData.logs"
            :key="index"
            :timestamp="log.time"
            :type="log.level === 'error' ? 'danger' : log.level === 'warn' ? 'warning' : 'primary'"
          >
            {{ log.message }}
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ResourcePoolTasks' })

const loading = ref(false)
const stats = ref({ pending: 0, running: 0, failed: 0 })

const searchForm = reactive({
  task_type: undefined as string | undefined,
  status: undefined as string | undefined,
  date_range: [] as string[]
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref([])
const detailDialogVisible = ref(false)
const detailData = ref<any>({})

const getTaskTypeTag = (type: string) => {
  const map: Record<string, any> = { allocate: 'primary', reclaim: 'warning', migrate: 'success', expand: 'info', health_check: '' }
  return map[type] || 'info'
}

const getTaskTypeText = (type: string) => {
  const map: Record<string, string> = {
    allocate: $t('clientsResourcePool.taskAllocate'),
    reclaim: $t('clientsResourcePool.taskReclaim'),
    migrate: $t('clientsResourcePool.taskMigrate'),
    expand: $t('clientsResourcePool.taskExpand'),
    health_check: $t('clientsResourcePool.taskHealthCheck')
  }
  return map[type] || $t('clientsResourcePool.unknown')
}

const getStatusTag = (status: string) => {
  const map: Record<string, any> = { pending: 'info', running: 'primary', completed: 'success', failed: 'danger', cancelled: 'warning' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: $t('clientsResourcePool.statusPending'),
    running: $t('clientsResourcePool.statusRunning'),
    completed: $t('clientsResourcePool.statusCompleted'),
    failed: $t('clientsResourcePool.statusFailed'),
    cancelled: $t('clientsResourcePool.statusCancelled')
  }
  return map[status] || $t('clientsResourcePool.unknown')
}

const fetchTasks = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.task_type) params.task_type = searchForm.task_type
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/resource-pool/tasks', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
    stats.value = { pending: data.pending || 0, running: data.running || 0, failed: data.failed || 0 }
  } catch (error) {
    ElMessage.error($t('clientsResourcePool.fetchTasksFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchTasks() }
const handleReset = () => { searchForm.task_type = undefined; searchForm.status = undefined; searchForm.date_range = []; handleSearch() }

const handleViewDetail = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/resource-pool/tasks/${row.id}` })
    detailData.value = data || row
    detailDialogVisible.value = true
  } catch (error) {
    detailData.value = row
    detailDialogVisible.value = true
  }
}

const handleRetry = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/resource-pool/tasks/${row.id}/retry`, showSuccessMessage: true })
    ElMessage.success($t('clientsResourcePool.taskResubmitted'))
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('clientsResourcePool.retryFailed'))
  }
}

const handleCancel = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/resource-pool/tasks/${row.id}/cancel`, showSuccessMessage: true })
    ElMessage.success($t('clientsResourcePool.taskCancelled'))
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('clientsResourcePool.cancelFailed'))
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchTasks() }
const handlePageChange = () => { fetchTasks() }

onMounted(() => { fetchTasks() })
</script>

<style scoped lang="scss">
.resource-pool-tasks-page { padding: 20px; }
.search-card { margin-bottom: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { .el-form-item { margin-bottom: 0; } }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.text-danger { color: #f56c6c; }
.detail-pre { margin: 0; font-size: 12px; white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow-y: auto; }
.task-logs { h4 { margin: 0 0 16px 0; font-size: 14px; color: #303133; } }
</style>
