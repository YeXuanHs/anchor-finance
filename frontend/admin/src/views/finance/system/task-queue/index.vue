<template>
  <div class="task-queue-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('taskQueue.title') }}</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleClearFailed">{{ $t('taskQueue.clearFailed') }}</el-button>
            <el-button type="primary" size="small" @click="handleRetryAll">{{ $t('taskQueue.retryAll') }}</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20" class="stat-cards">
        <el-col :span="4"><el-card shadow="hover" class="stat-card"><div class="stat-value">{{ stats.pending }}</div><div class="stat-label">{{ $t('taskQueue.pending') }}</div></el-card></el-col>
        <el-col :span="4"><el-card shadow="hover" class="stat-card processing"><div class="stat-value">{{ stats.processing }}</div><div class="stat-label">{{ $t('taskQueue.processing') }}</div></el-card></el-col>
        <el-col :span="4"><el-card shadow="hover" class="stat-card success"><div class="stat-value">{{ stats.completed }}</div><div class="stat-label">{{ $t('taskQueue.completed') }}</div></el-card></el-col>
        <el-col :span="4"><el-card shadow="hover" class="stat-card danger"><div class="stat-value">{{ stats.failed }}</div><div class="stat-label">{{ $t('taskQueue.failed') }}</div></el-card></el-col>
        <el-col :span="4"><el-card shadow="hover" class="stat-card"><div class="stat-value">{{ stats.delayed }}</div><div class="stat-label">{{ $t('taskQueue.delayed') }}</div></el-card></el-col>
        <el-col :span="4"><el-card shadow="hover" class="stat-card"><div class="stat-value">{{ stats.total }}</div><div class="stat-label">{{ $t('taskQueue.total') }}</div></el-card></el-col>
      </el-row>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('taskQueue.taskType')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('taskQueue.typeEmail')" value="email" />
            <el-option :label="$t('taskQueue.typeSms')" value="sms" />
            <el-option :label="$t('taskQueue.typeOrder')" value="order" />
            <el-option :label="$t('taskQueue.typeSync')" value="sync" />
            <el-option :label="$t('taskQueue.typeReport')" value="report" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('taskQueue.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('taskQueue.pending')" value="pending" />
            <el-option :label="$t('taskQueue.processing')" value="processing" />
            <el-option :label="$t('taskQueue.completed')" value="completed" />
            <el-option :label="$t('taskQueue.failed')" value="failed" />
            <el-option :label="$t('taskQueue.delayed')" value="delayed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" :label="$t('taskQueue.taskType')" width="120">
          <template #default="{ row }"><el-tag size="small">{{ $t(`taskQueue.typeMap.${row.type}`) || row.type }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('taskQueue.taskName')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('taskQueue.status')" width="100">
          <template #default="{ row }"><el-tag :type="statusTypeMap[row.status]" size="small">{{ $t(`taskQueue.${row.status}`) }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="attempts" :label="$t('taskQueue.attempts')" width="80" align="center" />
        <el-table-column prop="max_attempts" :label="$t('taskQueue.maxAttempts')" width="80" align="center" />
        <el-table-column prop="run_at" :label="$t('taskQueue.runAt')" width="180" />
        <el-table-column prop="completed_at" :label="$t('taskQueue.completedAt')" width="180" />
        <el-table-column :label="$t('taskQueue.operations')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'failed'" type="success" link @click="handleRetry(row)">{{ $t('taskQueue.retry') }}</el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('taskQueue.detail') }}</el-button>
            <el-popconfirm v-if="row.status !== 'processing'" :title="$t('taskQueue.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" :title="$t('taskQueue.taskDetail')" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item :label="$t('taskQueue.taskId')">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.taskType')">{{ $t(`taskQueue.typeMap.${detailData.type}`) || detailData.type }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.taskName')" :span="2">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.status')"><el-tag :type="statusTypeMap[detailData.status]" size="small">{{ $t(`taskQueue.${detailData.status}`) }}</el-tag></el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.attempts')">{{ detailData.attempts }} / {{ detailData.max_attempts }}</el-descriptions-item>
        <el-descriptions-item :label="$t('common.createdAt')">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.runAt')">{{ detailData.run_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.completedAt')" :span="2">{{ detailData.completed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.errorInfo')" :span="2" v-if="detailData.error"><pre class="error-text">{{ detailData.error }}</pre></el-descriptions-item>
        <el-descriptions-item :label="$t('taskQueue.taskData')" :span="2"><pre class="data-text">{{ JSON.stringify(detailData.payload, null, 2) }}</pre></el-descriptions-item>
      </el-descriptions>
      <template #footer><el-button @click="detailVisible = false">{{ $t('common.close') }}</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const statusTypeMap: Record<string, any> = { pending: 'info', processing: 'warning', completed: 'success', failed: 'danger', delayed: 'info' }

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const stats = reactive({ pending: 0, processing: 0, completed: 0, failed: 0, delayed: 0, total: 0 })
const searchForm = reactive({ type: '', status: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/system/task-queue', params: { page: pagination.page, page_size: pagination.page_size, ...searchForm } }); tableData.value = data.list || []; pagination.total = data.total || 0 } catch { ElMessage.error($t('taskQueue.fetchFailed')) } finally { loading.value = false }
}

const fetchStats = async () => { try { const data = await request.get({ url: '/api/admin/system/task-queue/stats' }); if (data) Object.assign(stats, data) } catch {} }

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.type = ''; searchForm.status = ''; handleSearch() }
const handleViewDetail = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleRetry = async (row: any) => { try { await request.post({ url: `/api/admin/system/task-queue/${row.id}/retry` }); ElMessage.success($t('taskQueue.retrySuccess')); fetchData(); fetchStats() } catch { ElMessage.error($t('common.operateFailed')) } }
const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/system/task-queue/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchData(); fetchStats() } catch { ElMessage.error($t('common.deleteFailed')) } }
const handleClearFailed = async () => { try { await request.post({ url: '/api/admin/system/task-queue/clear-failed' }); ElMessage.success($t('taskQueue.clearFailedSuccess')); fetchData(); fetchStats() } catch { ElMessage.error($t('common.operateFailed')) } }
const handleRetryAll = async () => { try { await request.post({ url: '/api/admin/system/task-queue/retry-all' }); ElMessage.success($t('taskQueue.retryAllSuccess')); fetchData(); fetchStats() } catch { ElMessage.error($t('common.operateFailed')) } }

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }
onMounted(() => { fetchData(); fetchStats() })
</script>

<style scoped lang="scss">
.task-queue-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
.search-form { margin-top: 20px; margin-bottom: 20px; }
.stat-cards { margin-bottom: 20px; }
.stat-card { text-align: center; padding: 8px 0; .stat-value { font-size: 24px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; } .stat-label { color: var(--el-text-color-secondary); font-size: 14px; } &.processing .stat-value { color: var(--el-color-warning); } &.success .stat-value { color: var(--el-color-success); } &.danger .stat-value { color: var(--el-color-danger); } }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.error-text, .data-text { margin: 0; white-space: pre-wrap; word-break: break-all; font-size: 12px; max-height: 200px; overflow-y: auto; }
</style>
