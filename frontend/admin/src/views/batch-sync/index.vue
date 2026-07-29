<template>
  <div class="batch-sync-page page-container">
    <div class="page-header">
      <h2>批量同步</h2>
    </div>

    <div class="sync-container">
      <el-card class="sync-card">
        <template #header>
          <div class="card-header">
            <span>同步设置</span>
          </div>
        </template>
        <el-form :model="syncForm" label-width="120px">
          <el-form-item label="同步类型">
            <el-select v-model="syncForm.type" placeholder="请选择同步类型">
              <el-option label="产品同步" value="product" />
              <el-option label="用户同步" value="user" />
              <el-option label="订单同步" value="order" />
            </el-select>
          </el-form-item>
          <el-form-item label="同步方向">
            <el-radio-group v-model="syncForm.direction">
              <el-radio label="import">导入</el-radio>
              <el-radio label="export">导出</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="数据源">
            <el-select v-model="syncForm.source" placeholder="请选择数据源">
              <el-option label="智简魔方" value="zjmf" />
              <el-option label="本地数据库" value="local" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSync" :loading="syncing">开始同步</el-button>
            <el-button @click="handleTest">测试连接</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card class="log-card">
        <template #header>
          <div class="card-header">
            <span>同步日志</span>
            <el-button type="text" @click="clearLogs">清空日志</el-button>
          </div>
        </template>
        <div class="log-content" ref="logContainer">
          <div v-for="(log, index) in logs" :key="index" :class="['log-item', log.type]">
            <span class="log-time">{{ log.time }}</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
          <div v-if="logs.length === 0" class="empty-log">暂无日志</div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const syncing = ref(false)
const logContainer = ref(null)

const syncForm = reactive({
  type: '',
  direction: 'import',
  source: ''
})

const logs = ref([])

const addLog = (message, type = 'info') => {
  const now = new Date()
  const time = now.toLocaleTimeString()
  logs.value.push({ time, message, type })
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const handleSync = async () => {
  if (!syncForm.type) {
    ElMessage.warning('请选择同步类型')
    return
  }

  syncing.value = true
  addLog('开始同步...', 'info')

  try {
    const { data } = await request.post('/admin/api/v1/batch-sync', syncForm)
    addLog('同步完成', 'success')
    ElMessage.success('同步完成')
  } catch (error) {
    addLog(`同步失败: ${error.message}`, 'error')
    ElMessage.error('同步失败')
  } finally {
    syncing.value = false
  }
}

const handleTest = async () => {
  try {
    await request.post('/admin/api/v1/batch-sync/test', syncForm)
    addLog('测试连接成功', 'success')
    ElMessage.success('连接测试成功')
  } catch (error) {
    addLog(`连接测试失败: ${error.message}`, 'error')
    ElMessage.error('连接测试失败')
  }
}

const clearLogs = () => {
  logs.value = []
}
</script>

<style scoped>
.batch-sync-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.sync-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.sync-card,
.log-card {
  height: fit-content;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-content {
  height: 400px;
  overflow-y: auto;
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
}

.log-item {
  margin-bottom: 5px;
  line-height: 1.5;
}

.log-time {
  color: #909399;
  margin-right: 10px;
}

.log-message {
  color: #606266;
}

.log-item.success .log-message {
  color: #67c23a;
}

.log-item.error .log-message {
  color: #f56c6c;
}

.log-item.warning .log-message {
  color: #e6a23c;
}

.empty-log {
  text-align: center;
  color: #909399;
  padding: 20px;
}
</style>
