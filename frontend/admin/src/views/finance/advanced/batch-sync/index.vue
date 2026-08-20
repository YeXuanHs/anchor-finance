<template>
  <div class="page-container">
    <art-card :title="$t('finance.batchSync.pageTitle')" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('finance.batchSync.taskName')" min-width="150" />
        <el-table-column prop="type" :label="$t('finance.batchSync.syncType')" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.batchSync.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'info'">
              {{ row.status === 0 ? $t('finance.batchSync.statusPending') : row.status === 1 ? $t('finance.batchSync.statusCompleted') : $t('finance.batchSync.statusRunning') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total" :label="$t('finance.batchSync.total')" width="100" />
        <el-table-column prop="synced" :label="$t('finance.batchSync.synced')" width="100" />
        <el-table-column prop="failed" :label="$t('finance.batchSync.failed')" width="100" />
        <el-table-column prop="created_at" :label="$t('finance.batchSync.createdAt')" width="180" />
        <el-table-column :label="$t('finance.batchSync.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleExecute(row)" :disabled="row.status === 2">{{ $t('finance.batchSync.execute') }}</el-button>
            <el-button size="small" @click="handleLogs(row)">{{ $t('finance.batchSync.logs') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <!-- 同步日志对话框 -->
    <el-dialog v-model="logsDialogVisible" :title="$t('finance.batchSync.syncLogs')" width="800px">
      <el-table :data="logsList" v-loading="logsLoading" stripe max-height="500">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="level" :label="$t('finance.batchSync.level')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.level === 'error' ? 'danger' : row.level === 'warn' ? 'warning' : 'info'" size="small">
              {{ row.level }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="$t('finance.batchSync.logContent')" min-width="400" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('finance.batchSync.time')" width="170" />
      </el-table>
      <template #footer>
        <el-button @click="logsDialogVisible = false">{{ $t('finance.batchSync.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])

// 日志对话框
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logsList = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/batch-sync' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleExecute = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/batch-sync/${row.id}/execute` })
    ElMessage.success($t('finance.batchSync.startExecute'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleLogs = async (row: any) => {
  logsDialogVisible.value = true
  logsLoading.value = true
  try {
    const res = await request.get({ url: `/api/admin/batch-sync/${row.id}/logs` })
    logsList.value = res || []
  } catch (error) {
    ElMessage.error($t('finance.batchSync.fetchLogsFailed'))
  } finally {
    logsLoading.value = false
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
