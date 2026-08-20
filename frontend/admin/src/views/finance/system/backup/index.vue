<template>
  <div class="backup-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('backup.title') }}</span>
          <el-button type="primary" @click="handleCreateBackup" :loading="createLoading">
            <el-icon><Download /></el-icon>
            {{ $t('backup.createBackup') }}
          </el-button>
        </div>
      </template>

      <!-- 备份统计 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total_count }}</div>
            <div class="stat-label">{{ $t('backup.totalBackups') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total_size }}</div>
            <div class="stat-label">{{ $t('backup.totalSize') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.last_backup_time || '-' }}</div>
            <div class="stat-label">{{ $t('backup.lastBackup') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.auto_backup_enabled ? $t('backup.enabled') : $t('backup.notEnabled') }}</div>
            <div class="stat-label">{{ $t('backup.autoBackup') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 备份任务状态 -->
      <div class="section" v-if="currentTask">
        <div class="section-header">
          <h3>{{ $t('backup.activeTask') }}</h3>
        </div>
        <el-card shadow="hover" class="task-card">
          <div class="task-info">
            <span>{{ $t('backup.taskId') }}: {{ currentTask.id }}</span>
            <el-tag :type="getTaskStatusType(currentTask.status)" size="small">
              {{ getTaskStatusLabel(currentTask.status) }}
            </el-tag>
          </div>
          <el-progress :percentage="currentTask.progress" :status="currentTask.status === 'failed' ? 'exception' : undefined" />
          <div class="task-actions">
            <el-button v-if="currentTask.status === 'running'" type="danger" size="small" @click="handleCancelTask(currentTask)">
              {{ $t('backup.cancelTask') }}
            </el-button>
          </div>
        </el-card>
      </div>

      <!-- 备份历史列表 -->
      <div class="section">
        <div class="section-header">
          <h3>{{ $t('backup.backupHistory') }}</h3>
        </div>
        <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" :label="$t('backup.id')" width="70" />
          <el-table-column prop="name" :label="$t('backup.backupFileName')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="size" :label="$t('backup.fileSize')" width="120" align="center">
            <template #default="{ row }">{{ formatSize(row.size) }}</template>
          </el-table-column>
          <el-table-column prop="type" :label="$t('backup.backupType')" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ typeMap[row.type as keyof typeof typeMap] || row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" :label="$t('backup.status')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" :label="$t('backup.createTime')" width="180" />
          <el-table-column prop="completed_at" :label="$t('backup.completedTime')" width="180">
            <template #default="{ row }">{{ row.completed_at || '-' }}</template>
          </el-table-column>
          <el-table-column prop="operator" :label="$t('backup.operator')" width="120" />
          <el-table-column :label="$t('backup.operations')" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleDownload(row)" :disabled="row.status !== 'completed'">
                {{ $t('backup.download') }}
              </el-button>
              <el-button type="warning" link @click="handleRestore(row)" :disabled="row.status !== 'completed'">
                {{ $t('backup.restore') }}
              </el-button>
              <el-popconfirm :title="$t('backup.confirmDelete')" @confirm="handleDelete(row)">
                <template #reference><el-button type="danger" link>{{ $t('backup.delete') }}</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
            :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange" @current-change="handlePageChange" />
        </div>
      </div>
    </el-card>

    <!-- 自动备份配置对话框 -->
    <el-dialog v-model="configDialogVisible" :title="$t('backup.autoBackupConfig')" width="500px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item :label="$t('backup.enableAutoBackup')">
          <el-switch v-model="configForm.enabled" />
        </el-form-item>
        <el-form-item :label="$t('backup.backupCycle')">
          <el-select v-model="configForm.schedule" style="width: 100%">
            <el-option :label="$t('backup.daily')" value="daily" />
            <el-option :label="$t('backup.weekly')" value="weekly" />
            <el-option :label="$t('backup.monthly')" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('backup.keepCount')">
          <el-input-number v-model="configForm.keep_count" :min="1" :max="30" />
          <span class="form-tip">{{ $t('backup.keepCountTip') }}</span>
        </el-form-item>
        <el-form-item :label="$t('backup.backupType')">
          <el-checkbox-group v-model="configForm.types">
            <el-checkbox value="database">{{ $t('backup.database') }}</el-checkbox>
            <el-checkbox value="files">{{ $t('backup.files') }}</el-checkbox>
            <el-checkbox value="config">{{ $t('backup.config') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">{{ $t('backup.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="configLoading">{{ $t('backup.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const typeMap = computed(() => ({
  database: $t('backup.database'),
  files: $t('backup.files'),
  config: $t('backup.config'),
  full: $t('backup.fullBackup')
}))
const statusMap = computed(() => ({
  completed: { label: $t('backup.completed'), type: 'success' },
  running: { label: $t('backup.running'), type: 'primary' },
  failed: { label: $t('backup.failed'), type: 'danger' },
  cancelled: { label: $t('backup.cancelled'), type: 'info' },
  pending: { label: $t('backup.pending'), type: 'warning' }
}))
const taskStatusMap = computed(() => ({
  running: { label: $t('backup.running'), type: 'primary' },
  failed: { label: $t('backup.failed'), type: 'danger' },
  cancelled: { label: $t('backup.cancelled'), type: 'info' }
}))

const getStatusType = (status: string) => (statusMap.value[status as keyof typeof statusMap.value]?.type || 'info') as any
const getStatusLabel = (status: string) => statusMap.value[status as keyof typeof statusMap.value]?.label || status
const getTaskStatusType = (status: string) => (taskStatusMap.value[status as keyof typeof taskStatusMap.value]?.type || 'info') as any
const getTaskStatusLabel = (status: string) => taskStatusMap.value[status as keyof typeof taskStatusMap.value]?.label || status

const formatSize = (kb: number) => {
  if (!kb) return '-'
  if (kb < 1024) return kb + ' KB'
  return (kb / 1024).toFixed(2) + ' MB'
}

const loading = ref(false)
const createLoading = ref(false)
const configLoading = ref(false)
const configDialogVisible = ref(false)

const tableData = ref<any[]>([])
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const stats = reactive({ total_count: 0, total_size: '0 MB', last_backup_time: '', auto_backup_enabled: false })
const currentTask = ref<any>(null)

const configForm = reactive({ enabled: false, schedule: 'daily', keep_count: 7, types: ['database'] })

const fetchData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.page, page_size: pagination.page_size }
    const res = await request.get({ url: '/api/admin/backups', params })
    tableData.value = res?.list || res?.data || res || []
    pagination.total = res?.total || 0
  } catch { ElMessage.error($t('backup.fetchListFailed')) } finally { loading.value = false }
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/backup/stats' })
    if (data) Object.assign(stats, data)
  } catch { /* ignore */ }
}

const fetchCurrentTask = async () => {
  try {
    const data = await request.get({ url: '/api/admin/backup/task/current' })
    currentTask.value = data || null
  } catch { /* ignore */ }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleCreateBackup = async () => {
  try {
    await ElMessageBox.confirm($t('backup.confirmCreate'), $t('backup.createBackupTitle'))
    createLoading.value = true
    await request.post({ url: '/api/admin/backups', showSuccessMessage: true })
    fetchData(); fetchStats(); fetchCurrentTask()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('backup.createFailed')) } finally { createLoading.value = false }
}

const handleCancelTask = async (task: any) => {
  try {
    await ElMessageBox.confirm($t('backup.confirmCancelTask'), $t('backup.cancelTaskTitle'))
    await request.post({ url: `/api/admin/backups/task/${task.id}/cancel` })
    ElMessage.success($t('backup.taskCancelled')); fetchCurrentTask()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('backup.cancelFailed')) }
}

const handleDownload = async (row: any) => {
  try {
    const res = await request.get({ url: `/api/admin/backups/${row.id}/download` })
    if (res?.url) window.open(res.url)
  } catch { ElMessage.error($t('backup.downloadFailed')) }
}

const handleRestore = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('backup.confirmRestore', { name: row.name }), $t('backup.restoreBackupTitle'), { type: 'warning' })
    await request.post({ url: `/api/admin/backups/${row.id}/restore` })
    ElMessage.success($t('backup.restoreSubmitted'))
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('backup.restoreFailed')) }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/backups/${row.id}` })
    ElMessage.success($t('backup.deleteSuccess')); fetchData(); fetchStats()
  } catch { ElMessage.error($t('backup.deleteFailed')) }
}

const handleSaveConfig = async () => {
  configLoading.value = true
  try {
    await request.put({ url: '/api/admin/backups/config', data: configForm, showSuccessMessage: true })
    configDialogVisible.value = false; fetchStats()
  } catch { ElMessage.error($t('backup.saveFailed')) } finally { configLoading.value = false }
}

onMounted(() => { fetchData(); fetchStats(); fetchCurrentTask() })
</script>

<style scoped lang="scss">
.backup-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.stat-section { margin-bottom: 24px; }
.stat-card { text-align: center; padding: 8px 0; }
.stat-value { font-size: 20px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; }
.stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
.section { margin-top: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { margin: 0; font-size: 16px; font-weight: 600; } }
.task-card { margin-bottom: 16px; }
.task-info { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.task-actions { margin-top: 12px; text-align: right; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.form-tip { margin-left: 12px; color: var(--el-text-color-secondary); font-size: 13px; }
</style>
