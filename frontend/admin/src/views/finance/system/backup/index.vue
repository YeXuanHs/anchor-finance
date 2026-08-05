<template>
  <div class="backup-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>备份管理</span>
          <el-button type="primary" @click="handleCreateBackup" :loading="createLoading">
            <el-icon><Download /></el-icon>
            创建备份
          </el-button>
        </div>
      </template>

      <!-- 备份统计 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total_count }}</div>
            <div class="stat-label">备份总数</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.total_size }}</div>
            <div class="stat-label">总大小</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.last_backup_time || '-' }}</div>
            <div class="stat-label">上次备份</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ stats.auto_backup_enabled ? '已启用' : '未启用' }}</div>
            <div class="stat-label">自动备份</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 备份任务状态 -->
      <div class="section" v-if="currentTask">
        <div class="section-header">
          <h3>进行中的备份任务</h3>
        </div>
        <el-card shadow="hover" class="task-card">
          <div class="task-info">
            <span>任务ID: {{ currentTask.id }}</span>
            <el-tag :type="getTaskStatusType(currentTask.status)" size="small">
              {{ getTaskStatusLabel(currentTask.status) }}
            </el-tag>
          </div>
          <el-progress :percentage="currentTask.progress" :status="currentTask.status === 'failed' ? 'exception' : undefined" />
          <div class="task-actions">
            <el-button v-if="currentTask.status === 'running'" type="danger" size="small" @click="handleCancelTask(currentTask)">
              取消任务
            </el-button>
          </div>
        </el-card>
      </div>

      <!-- 备份历史列表 -->
      <div class="section">
        <div class="section-header">
          <h3>备份历史</h3>
        </div>
        <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="name" label="备份文件名" min-width="200" show-overflow-tooltip />
          <el-table-column prop="size" label="文件大小" width="120" align="center">
            <template #default="{ row }">{{ formatSize(row.size) }}</template>
          </el-table-column>
          <el-table-column prop="type" label="备份类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column prop="completed_at" label="完成时间" width="180">
            <template #default="{ row }">{{ row.completed_at || '-' }}</template>
          </el-table-column>
          <el-table-column prop="operator" label="操作人" width="120" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleDownload(row)" :disabled="row.status !== 'completed'">
                下载
              </el-button>
              <el-button type="warning" link @click="handleRestore(row)" :disabled="row.status !== 'completed'">
                恢复
              </el-button>
              <el-popconfirm title="确定删除该备份文件吗？删除后不可恢复。" @confirm="handleDelete(row)">
                <template #reference><el-button type="danger" link>删除</el-button></template>
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
    <el-dialog v-model="configDialogVisible" title="自动备份配置" width="500px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item label="启用自动备份">
          <el-switch v-model="configForm.enabled" />
        </el-form-item>
        <el-form-item label="备份周期">
          <el-select v-model="configForm.schedule" style="width: 100%">
            <el-option label="每天" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item label="保留份数">
          <el-input-number v-model="configForm.keep_count" :min="1" :max="30" />
          <span class="form-tip">超过此数量的旧备份将自动删除</span>
        </el-form-item>
        <el-form-item label="备份类型">
          <el-checkbox-group v-model="configForm.types">
            <el-checkbox value="database">数据库</el-checkbox>
            <el-checkbox value="files">文件</el-checkbox>
            <el-checkbox value="config">配置</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="configLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const typeMap: Record<string, string> = { database: '数据库', files: '文件', config: '配置', full: '完整备份' }
const statusMap: Record<string, { label: string; type: string }> = {
  completed: { label: '已完成', type: 'success' },
  running: { label: '进行中', type: 'primary' },
  failed: { label: '失败', type: 'danger' },
  cancelled: { label: '已取消', type: 'info' },
  pending: { label: '等待中', type: 'warning' }
}
const taskStatusMap: Record<string, { label: string; type: string }> = {
  running: { label: '进行中', type: 'primary' },
  failed: { label: '失败', type: 'danger' },
  cancelled: { label: '已取消', type: 'info' }
}

const getStatusType = (status: string) => (statusMap[status]?.type || 'info') as any
const getStatusLabel = (status: string) => statusMap[status]?.label || status
const getTaskStatusType = (status: string) => (taskStatusMap[status]?.type || 'info') as any
const getTaskStatusLabel = (status: string) => taskStatusMap[status]?.label || status

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
    const res = await request.get({ url: '/api/admin/backup', params })
    tableData.value = res?.list || res?.data || res || []
    pagination.total = res?.total || 0
  } catch { ElMessage.error('获取备份列表失败') } finally { loading.value = false }
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
    await ElMessageBox.confirm('确定创建数据库备份吗？此操作可能需要一些时间。', '创建备份')
    createLoading.value = true
    await request.post({ url: '/api/admin/backup', showSuccessMessage: true })
    fetchData(); fetchStats(); fetchCurrentTask()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('创建备份失败') } finally { createLoading.value = false }
}

const handleCancelTask = async (task: any) => {
  try {
    await ElMessageBox.confirm('确定取消当前备份任务吗？', '取消任务')
    await request.post({ url: `/api/admin/backup/task/${task.id}/cancel` })
    ElMessage.success('任务已取消'); fetchCurrentTask()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('取消失败') }
}

const handleDownload = async (row: any) => {
  try {
    const res = await request.get({ url: `/api/admin/backup/${row.id}/download` })
    if (res?.url) window.open(res.url)
  } catch { ElMessage.error('下载失败') }
}

const handleRestore = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定恢复备份 "${row.name}" 吗？此操作将覆盖当前数据，请谨慎操作！`, '恢复备份', { type: 'warning' })
    await request.post({ url: `/api/admin/backup/${row.id}/restore` })
    ElMessage.success('恢复任务已提交')
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('恢复失败') }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/backup/${row.id}` })
    ElMessage.success('删除成功'); fetchData(); fetchStats()
  } catch { ElMessage.error('删除失败') }
}

const handleSaveConfig = async () => {
  configLoading.value = true
  try {
    await request.put({ url: '/api/admin/backup/config', data: configForm, showSuccessMessage: true })
    configDialogVisible.value = false; fetchStats()
  } catch { ElMessage.error('保存失败') } finally { configLoading.value = false }
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
