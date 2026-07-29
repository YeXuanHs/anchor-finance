<template>
  <div class="cron-url-page page-container">
    <div class="page-header">
      <h2>URL定时任务</h2>
      <el-button type="primary" @click="handleRunCron">立即执行</el-button>
    </div>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>定时任务配置</span>
        </div>
      </template>
      <el-form :model="cronForm" label-width="120px">
        <el-form-item label="任务URL">
          <el-input v-model="cronForm.url" placeholder="请输入定时任务URL" disabled>
            <template #append>
              <el-button @click="handleCopyUrl">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="执行间隔">
          <el-select v-model="cronForm.interval" placeholder="请选择执行间隔">
            <el-option label="每分钟" value="1" />
            <el-option label="每5分钟" value="5" />
            <el-option label="每10分钟" value="10" />
            <el-option label="每30分钟" value="30" />
            <el-option label="每小时" value="60" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="cronForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>执行日志</span>
          <el-button type="text" @click="fetchLogs">刷新</el-button>
        </div>
      </template>
      <el-table :data="logs" border fit>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="executed_at" label="执行时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="执行结果" />
        <el-table-column prop="duration" label="耗时(秒)" width="100" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const cronForm = reactive({
  url: '',
  interval: '5',
  enabled: true
})

const logs = ref([])

const fetchConfig = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/cron-url/config')
    if (data.data) {
      Object.assign(cronForm, data.data)
    }
  } catch (error) {
    ElMessage.error('获取配置失败')
  }
}

const fetchLogs = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/cron-url/logs')
    logs.value = data.data?.list || data.data || []
  } catch (error) {
    ElMessage.error('获取日志失败')
  }
}

const handleRunCron = async () => {
  try {
    await request.post('/admin/api/v1/cron-url/run')
    ElMessage.success('定时任务执行成功')
    fetchLogs()
  } catch (error) {
    ElMessage.error('定时任务执行失败')
  }
}

const handleCopyUrl = () => {
  navigator.clipboard.writeText(cronForm.url)
  ElMessage.success('URL已复制到剪贴板')
}

onMounted(() => {
  fetchConfig()
  fetchLogs()
})
</script>

<style scoped>
.cron-url-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
