<template>
  <div class="batch-sync-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>批量同步</h3>
        <div>
          <el-button type="primary" @click="startSync">
            <el-icon><Refresh /></el-icon>
            开始同步
          </el-button>
          <el-button @click="selectSync">
            <el-icon><Select /></el-icon>
            选择同步
          </el-button>
        </div>
      </div>

      <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
        批量同步将从上游服务器同步服务器信息、状态和资源使用情况。同步过程可能需要一些时间，请耐心等待。
      </el-alert>

      <el-table :data="syncTasks" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="upstream_name" label="上游" />
        <el-table-column prop="server_count" label="服务器数" width="100" />
        <el-table-column prop="last_sync" label="最后同步" width="180" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getSyncStatusType(row.status)" size="small">
              {{ getSyncStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="200">
          <template #default="{ row }">
            <el-progress :percentage="row.progress" :status="row.progress === 100 ? 'success' : undefined" />
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="art-card" style="margin-top: 20px;">
      <h4>同步日志</h4>
      <el-timeline>
        <el-timeline-item v-for="log in syncLogs" :key="log.id" :timestamp="log.time" placement="top" :type="log.type === 'success' ? 'success' : log.type === 'error' ? 'danger' : 'primary'">
          <p>{{ log.message }}</p>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Select } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const syncTasks = ref([])
const selectedTasks = ref<any[]>([])
const syncLogs = ref<any[]>([])

const getSyncStatusType = (status: string) => {
  const map: Record<string, string> = { idle: 'info', syncing: 'warning', success: 'success', error: 'danger' }
  return map[status] || 'info'
}

const getSyncStatusText = (status: string) => {
  const map: Record<string, string> = { idle: '空闲', syncing: '同步中', success: '成功', error: '失败' }
  return map[status] || status
}

const fetchSyncTasks = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/batch-sync')
    if (data?.data) {
      syncTasks.value = data.data.tasks || data.data
      syncLogs.value = data.data.logs || []
    }
  } catch {
    ElMessage.error('获取同步任务失败')
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection: any[]) => { selectedTasks.value = selection }

const startSync = async () => {
  if (!selectedTasks.value.length) {
    ElMessage.warning('请先选择要同步的任务')
    return
  }
  try {
    const ids = selectedTasks.value.map((t: any) => t.id)
    await request.post('/admin/batch-sync/execute', { ids })
    ElMessage.success('同步任务已启动')
    fetchSyncTasks()
  } catch {
    ElMessage.error('启动同步失败')
  }
}

const selectSync = async () => {
  if (!selectedTasks.value.length) {
    ElMessage.warning('请先选择要同步的任务')
    return
  }
  for (const task of selectedTasks.value) {
    try {
      await request.post(`/admin/batch-sync/${task.id}/execute`)
    } catch {}
  }
  ElMessage.success('选中的同步任务已全部启动')
  fetchSyncTasks()
}

onMounted(() => {
  fetchSyncTasks()
})
</script>

<style scoped lang="scss">
.batch-sync-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  h4 { margin: 0 0 16px; font-size: 15px; font-weight: 600; }
}
</style>
